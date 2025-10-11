# MongoDB User Setup

## Creating the Initial Admin User

Before using the authentication endpoints, you need to create at least one user in your MongoDB database.

### Using MongoDB Compass or MongoDB Shell

Connect to your MongoDB instance and run the following command in the `trimanaDB` database:

```javascript
db.users.insertOne({
  "_id": "admin-user-1",
  "username": "admin",
  "password": "your-secure-password"
})
```

### Security Note

**IMPORTANT:** In production, passwords should be hashed before storing them in the database. 

Consider using bcrypt to hash passwords:

```go
import "golang.org/x/crypto/bcrypt"

// Hashing a password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verifying a password
err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(providedPassword))
```

### Updating auth.go to use bcrypt (Optional but Recommended)

1. Add bcrypt to your go.mod:
```bash
go get golang.org/x/crypto/bcrypt
```

2. Update `payroll/auth.go` to hash and verify passwords:

```go
package payroll

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"` // This should be hashed
}

var userCollection = "users"

func getUserCollection() *mongo.Collection {
	return MongoClient.Database(db).Collection(userCollection)
}

func AuthenticateUser(username, password string) (*User, error) {
	collection := getUserCollection()
	var user User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return nil, err
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

// Helper function to hash passwords (use this when creating users)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
```

3. When creating a new user, hash the password first:

```javascript
// Generate hashed password using bcrypt (e.g., using an online tool or script)
// For example, the bcrypt hash of "admin" is:
// $2a$10$N9qo8uLOickgx2ZMRZoMye4UI.JqjXf5fKkCXsRE3VJW7gQGJQqKm

db.users.insertOne({
  "_id": "admin-user-1",
  "username": "admin",
  "password": "$2a$10$N9qo8uLOickgx2ZMRZoMye4UI.JqjXf5fKkCXsRE3VJW7gQGJQqKm"
})
```

## Testing Authentication

Once you have created a user, you can test the login endpoint:

```bash
curl -X POST https://your-api-url/prod/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

Expected response:
```json
{
  "username": "admin",
  "token": "generated-token-string"
}
```

