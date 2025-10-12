package payroll

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type EmployeeStats struct {
	EmployeeID       string  `json:"employeeId"`
	EmployeeName     string  `json:"employeeName"`
	TotalHours       float64 `json:"totalHours"`
	OvertimeHours    float64 `json:"overtimeHours"`
	TotalPay         float64 `json:"totalPay"`
	DaysWorked       int     `json:"daysWorked"`
	AveragePerDay    float64 `json:"averagePerDay"`
	TotalBreakTime   float64 `json:"totalBreakTime"`
	AverageBreakTime float64 `json:"averageBreakTime"`
	BreakCount       int     `json:"breakCount"`
}

type DailyStats struct {
	Date       time.Time `json:"date"`
	TotalHours float64   `json:"totalHours"`
	Employees  int       `json:"employees"`
}

type OverallStats struct {
	TotalEmployees       int             `json:"totalEmployees"`
	ActiveEmployees      int             `json:"activeEmployees"`
	TotalHoursThisPeriod float64         `json:"totalHoursThisPeriod"`
	TotalPayrollCost     float64         `json:"totalPayrollCost"`
	TotalOvertimeHours   float64         `json:"totalOvertimeHours"`
	AverageHoursPerEmp   float64         `json:"averageHoursPerEmp"`
	TotalBreakTime       float64         `json:"totalBreakTime"`
	AverageBreakPerEmp   float64         `json:"averageBreakPerEmp"`
	TopPerformers        []EmployeeStats `json:"topPerformers"`
	LeastActive          []EmployeeStats `json:"leastActive"`
	LongestBreaks        []EmployeeStats `json:"longestBreaks"`
	NoBreakEmployees     []EmployeeStats `json:"noBreakEmployees"`
	DailyTrend           []DailyStats    `json:"dailyTrend"`
	EmployeeBreakdown    []EmployeeStats `json:"employeeBreakdown"`
}

func GetPayrollStats(startDate, endDate time.Time) (*OverallStats, error) {
	employees, err := GetEmployees()
	if err != nil {
		return nil, err
	}

	stats := &OverallStats{
		TotalEmployees:    len(employees),
		TopPerformers:     []EmployeeStats{},
		LeastActive:       []EmployeeStats{},
		LongestBreaks:     []EmployeeStats{},
		NoBreakEmployees:  []EmployeeStats{},
		DailyTrend:        []DailyStats{},
		EmployeeBreakdown: []EmployeeStats{},
	}

	var allEmployeeStats []EmployeeStats
	dailyHoursMap := make(map[string]*DailyStats)
	activeEmployees := 0

	for _, employee := range employees {
		timeCards, err := GetTimeCardsByPayPeriod(employee.EmployeeID, startDate, endDate)
		if err != nil {
			continue
		}

		totalHours := 0.0
		overtimeHours := 0.0
		totalBreakTime := 0.0
		breakCount := 0
		daysWorked := len(timeCards)

		for _, timeCard := range timeCards {
			dayHours := 0.0

			// Calculate break time between blocks
			for i := 0; i < len(timeCard.TimeBlocks); i++ {
				block := timeCard.TimeBlocks[i]
				if !block.CheckOut.IsZero() {
					hoursForBlock := block.CheckOut.Sub(block.CheckIn).Hours()
					totalHours += hoursForBlock
					dayHours += hoursForBlock

					// Track daily stats
					dateKey := timeCard.CreatedAt.Format("2006-01-02")
					if dailyStats, exists := dailyHoursMap[dateKey]; exists {
						dailyStats.TotalHours += hoursForBlock
					} else {
						dailyHoursMap[dateKey] = &DailyStats{
							Date:       timeCard.CreatedAt,
							TotalHours: hoursForBlock,
							Employees:  1,
						}
					}
				}

				// Calculate break time (gap between current block's checkout and next block's checkin)
				if i < len(timeCard.TimeBlocks)-1 && !block.CheckOut.IsZero() {
					nextBlock := timeCard.TimeBlocks[i+1]
					if !nextBlock.CheckIn.IsZero() {
						breakDuration := nextBlock.CheckIn.Sub(block.CheckOut).Hours()
						// Only count as break if less than 4 hours (otherwise might be end of shift)
						if breakDuration > 0 && breakDuration < 4 {
							totalBreakTime += breakDuration
							breakCount++
						}
					}
				}
			}

			// Calculate overtime for this day
			if dayHours > 8.0 {
				dayOvertime := dayHours - 8.0
				overtimeHours += dayOvertime
			}
		}

		if totalHours > 0 {
			activeEmployees++

			standardHours := totalHours - overtimeHours
			standardPay := standardHours * employee.HourlyRate
			overtimePay := overtimeHours * employee.HourlyRate * 1.5
			totalPay := standardPay + overtimePay

			averageBreakTime := 0.0
			if breakCount > 0 {
				averageBreakTime = totalBreakTime / float64(breakCount)
			}

			empStats := EmployeeStats{
				EmployeeID:       employee.EmployeeID,
				EmployeeName:     employee.Name,
				TotalHours:       totalHours,
				OvertimeHours:    overtimeHours,
				TotalPay:         totalPay,
				DaysWorked:       daysWorked,
				AveragePerDay:    totalHours / float64(daysWorked),
				TotalBreakTime:   totalBreakTime,
				AverageBreakTime: averageBreakTime,
				BreakCount:       breakCount,
			}

			allEmployeeStats = append(allEmployeeStats, empStats)
			stats.TotalHoursThisPeriod += totalHours
			stats.TotalOvertimeHours += overtimeHours
			stats.TotalPayrollCost += totalPay
			stats.TotalBreakTime += totalBreakTime
		}
	}

	stats.ActiveEmployees = activeEmployees
	if activeEmployees > 0 {
		stats.AverageHoursPerEmp = stats.TotalHoursThisPeriod / float64(activeEmployees)
		stats.AverageBreakPerEmp = stats.TotalBreakTime / float64(activeEmployees)
	}

	// Sort for top performers (most hours)
	stats.EmployeeBreakdown = allEmployeeStats

	// Get top 5 performers
	topPerformers := make([]EmployeeStats, len(allEmployeeStats))
	copy(topPerformers, allEmployeeStats)
	for i := 0; i < len(topPerformers)-1; i++ {
		for j := i + 1; j < len(topPerformers); j++ {
			if topPerformers[j].TotalHours > topPerformers[i].TotalHours {
				topPerformers[i], topPerformers[j] = topPerformers[j], topPerformers[i]
			}
		}
	}
	if len(topPerformers) > 5 {
		stats.TopPerformers = topPerformers[:5]
	} else {
		stats.TopPerformers = topPerformers
	}

	// Get least active (bottom 5)
	leastActive := make([]EmployeeStats, len(allEmployeeStats))
	copy(leastActive, allEmployeeStats)
	for i := 0; i < len(leastActive)-1; i++ {
		for j := i + 1; j < len(leastActive); j++ {
			if leastActive[j].TotalHours < leastActive[i].TotalHours {
				leastActive[i], leastActive[j] = leastActive[j], leastActive[i]
			}
		}
	}
	if len(leastActive) > 5 {
		stats.LeastActive = leastActive[:5]
	} else {
		stats.LeastActive = leastActive
	}

	// Get employees with longest breaks (top 5)
	longestBreaks := make([]EmployeeStats, len(allEmployeeStats))
	copy(longestBreaks, allEmployeeStats)
	for i := 0; i < len(longestBreaks)-1; i++ {
		for j := i + 1; j < len(longestBreaks); j++ {
			if longestBreaks[j].TotalBreakTime > longestBreaks[i].TotalBreakTime {
				longestBreaks[i], longestBreaks[j] = longestBreaks[j], longestBreaks[i]
			}
		}
	}
	if len(longestBreaks) > 5 {
		stats.LongestBreaks = longestBreaks[:5]
	} else {
		stats.LongestBreaks = longestBreaks
	}

	// Get employees with no breaks (compliance concern)
	var noBreakEmployees []EmployeeStats
	for _, emp := range allEmployeeStats {
		if emp.BreakCount == 0 && emp.TotalHours > 5.0 {
			noBreakEmployees = append(noBreakEmployees, emp)
		}
	}
	stats.NoBreakEmployees = noBreakEmployees

	// Convert daily stats map to sorted array
	for _, dailyStat := range dailyHoursMap {
		stats.DailyTrend = append(stats.DailyTrend, *dailyStat)
	}

	// Sort daily trend by date
	for i := 0; i < len(stats.DailyTrend)-1; i++ {
		for j := i + 1; j < len(stats.DailyTrend); j++ {
			if stats.DailyTrend[j].Date.Before(stats.DailyTrend[i].Date) {
				stats.DailyTrend[i], stats.DailyTrend[j] = stats.DailyTrend[j], stats.DailyTrend[i]
			}
		}
	}

	return stats, nil
}

func GetTodayStats() (int, float64, error) {
	dateNow := time.Now().In(TimeZone)
	dateNow = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, TimeZone)

	collection := getTimeCardCollection()
	cursor, err := collection.Find(context.Background(), bson.M{"createdAt": dateNow})
	if err != nil {
		return 0, 0, err
	}

	var timeCards []TimeCard
	if err = cursor.All(context.Background(), &timeCards); err != nil {
		return 0, 0, err
	}

	totalHours := 0.0
	for _, tc := range timeCards {
		for _, block := range tc.TimeBlocks {
			if !block.CheckOut.IsZero() {
				totalHours += block.CheckOut.Sub(block.CheckIn).Hours()
			}
		}
	}

	return len(timeCards), totalHours, nil
}
