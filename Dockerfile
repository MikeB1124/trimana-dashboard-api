FROM public.ecr.aws/lambda/go:latest


ENV AWS_ACCESS_KEY_ID=
ENV AWS_SECRET_ACCESS_KEY=
ENV AWS_SESSION_TOKEN=
# Copy function code
COPY main ${LAMBDA_TASK_ROOT}

# RUN chmod 755 ${LAMBDA_TASK_ROOT}/bootstrap

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "bootstrap" ]
