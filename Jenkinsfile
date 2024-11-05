pipeline {
    agent any
    
    tools {
        go 'go1.21.2'
    }
    
    environment {
        DB_HOST = credentials('DB_IP')
        DB_NAME = credentials('DB_NAME_elrek-system_DEV')
        DB_PASSWORD = credentials('DB_PASSWD')
        DB_PORT = credentials('DB_PORT')
        DB_USERNAME = credentials('DB_USERNAME')
        FRONTEND_URL = credentials('FRONTEND_URL_elrek-system')
        BACKEND_URL = credentials('BACKEND_URL_elrek-system')
        DOMAIN = credentials('DOMAIN_elrek-system')
        REPO = credentials('GITHUB_REPO_elrek-system')
    }
    
    stages {
        stage('Checkout') {
            steps {
                // Pull the latest code from the Git repository
                git branch: 'dev', url: REPO
            }
        }
        stage('Build') {
            steps {
                // Build the Go application
                sh 'go build -o elrek-system_GO'
            }
        }
        stage('Test') {
            steps {
                // Run Go unit tests
                sh 'go test -v ./tests'
            }
        }
    }
    post {
        always {
            cleanWs()
        }
    }
}