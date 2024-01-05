A **Clean Architecture** based web application that implements **CRUD** operations through **RESTful APIs**.
Docker containerised **PostgresQL** databse used for data persistence.
**Uber Zap** is used for logging. **Go Chi** Router is used for routing.

For running the application using Docker-Compose file, it is required to create an .env.development.local file.
The contents of the file are as follows:

-----------------------------[START of env file]--------------------------------  

POSTGRES_PASSWORD=secret
POSTGRES_DB=testdb  
POSTGRES_USER=nishant  

LIQUIBASE_COMMAND_URL="jdbc:postgresql://postgres:5432/testdb?user=nishant&password=secret"  
LIQUIBASE_CLASSPATH=/liquibase/changelog  
LIQUIBASE_COMMAND_CHANGELOG_FILE=db.changelog-root.xml  
LIQUIBASE_LOG_LEVEL=INFO  

Golang Server  
APP_ENV=development  
OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4318  

Jaeger Tracing  
COLLECTOR_OTLP_ENABLED=true  
  
--------------------------------[END of env file]--------------------------------  
