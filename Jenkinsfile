def dockerImage
def version

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
        SSH_HOST = credentials('HOST_SSH_elrek-system')
        SCRIPTS_HOST = credentials('HOST-SCRIPT-FOLDER_elrek-system')
        SSH_DB = credentials('DB_SSH')
        SCRIPTS_DB = credentials('DB-SCRIPT-FOLDER')
    }
    
    stages {
        // MARK: Checkout
        stage('Checkout DEV') {
            when {
                branch 'dev'
            }
            steps {
                git branch: 'dev', url: REPO
            }
        }
        stage('Checkout MAIN') {
            when {
                branch 'main'
            }
            steps {
                git branch: 'main', url: REPO
            }
        }

        // MARK: Read Version
        stage('Read Version') {
            steps {
                script {
                    version = readFile('version').trim()
                    echo "Building version ${version}"
                }
            }
        }

        // MARK: Build
        stage('Build') {
            steps {
                sh 'go build -o elrek-system_GO'
            }
        }

        // MARK: Test
        stage('Test') {
            steps {
                sh 'go test -v ./tests'
            }
        }

        // MARK: Build Docker image
        stage('Build Docker image') {
            steps {
                script {
                    dockerImage = docker.build('sc4n1a471/elrek-system_go')
                }
            }
        }

        // MARK: Push Docker image
        stage('Push production docker image') {
            when {
                branch 'main'
            }
            steps {
                script {
                    docker.withRegistry('https://registry.hub.docker.com', 'DOCKER_HUB') {
                        dockerImage.push("latest")
                        dockerImage.push("${version}")
                    }
                }
            }
        }
        stage('Push development docker image') {
            when {
                branch 'dev'
            }
            steps {
                script {
                    docker.withRegistry('https://registry.hub.docker.com', 'DOCKER_HUB') {
                        dockerImage.push("latest-dev")
                        dockerImage.push("${version}-dev")
                    }
                }
            }
        }

        // MARK: Backup DB
        stage('Backup DB') {
            when {
                anyOf {
                    branch 'main'
                    branch 'dev'
                }
            }

            steps {
                script {
                    echo "Backing up DB"

                    sh '''
                    ssh -tt $SSH_DB << EOF
                    cd $SCRIPTS_DB
                    ./elrek-backup.sh
                    exit
                    EOF'''
                }
            }
        }

        // MARK: Deploy
        stage('Deploy development') {
            when {
                branch 'dev'
            }

            steps {
                script {
                    echo "Deploying version ${version} to DEV"

                    sh """
                    ssh -tt \$SSH_HOST << EOF
                    cd \$SCRIPTS_HOST
                    ./redeploy-go.py '{"version": "${version}-dev", "env": "dev"}'
                    exit
                    EOF"""
                }
            }
        }
        stage('Deploy production') {
            when {
                branch 'main'
            }

            steps {
                script {
                    echo "Deploying version ${version} to PROD"

                    sh """
                    ssh -tt \$SSH_HOST << EOF
                    cd \$SCRIPTS_HOST
                    ./redeploy-go.py '{"version": "${version}", "env": "prod"}'
                    exit
                    EOF"""
                }
            }
        }
    }
    post {
        always {
            cleanWs()
        }
    }
}