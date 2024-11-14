def dockerImage
def version

pipeline {
    agent any
    
    tools {
        go 'go1.21.2'
    }
    
    environment {
        DB_HOST = credentials('DB_IP')

        DB_NAME = credentials('DB_NAME_elrek-system_DEV')           // For go tests
        DB_NAME_DEV = credentials('DB_NAME_elrek-system_DEV')
        DB_NAME_PROD = credentials('DB_NAME_elrek-system_PROD')

        DB_PASSWORD = credentials('DB_PASSWD')
        DB_PORT = credentials('DB_PORT')
        DB_USERNAME = credentials('DB_USERNAME')

        FRONTEND_URL = credentials('FRONTEND_URL_elrek-system_DEV') // For go tests
        BACKEND_URL = credentials('BACKEND_URL_elrek-system_DEV')   // For go tests
        FRONTEND_URL_DEV = credentials('FRONTEND_URL_elrek-system_DEV')
        BACKEND_URL_DEV = credentials('BACKEND_URL_elrek-system_DEV')
        FRONTEND_URL_PROD = credentials('FRONTEND_URL_elrek-system_PROD')
        BACKEND_URL_PROD = credentials('BACKEND_URL_elrek-system_PROD')

        DOMAIN = credentials('DOMAIN_elrek-system')
        REPO = credentials('GITHUB_REPO_elrek-system')
        SSH_HOST = credentials('HOST_SSH_elrek-system')
        SCRIPTS_HOST = credentials('HOST-SCRIPT-FOLDER_elrek-system')
        SSH_DB = credentials('DB_SSH')
        SCRIPTS_DB = credentials('DB-SCRIPT-FOLDER')

        GRAYLOG_HOST_DEV = credentials('GRAYLOG_HOST_elrek-system_DEV')
        GRAYLOG_HOST_PROD = credentials('GRAYLOG_HOST_elrek-system_PROD')
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
            when {
                anyOf {
                    branch 'main'
                    branch 'dev'
                }
            }
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

        // MARK: Test and Build Docker image
        stage('Test and Build Docker image') {
            parallel {
                stage('Test') {
                    steps {
                        sh 'go test -v ./tests'
                    }
                }
                stage('Build Docker image') {
                    steps {
                        script {
                            dockerImage = docker.build('sc4n1a471/elrek-system_go')
                        }
                    }
                }
            }
        }

        // MARK: Push and Backup
        stage('Push and Backup') {
            parallel {
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
                    ssh -tt $SSH_HOST << EOF
                    if [ "\$(docker ps -a -q -f name=elrek-system_go_dev)" ]; then
                        docker rm -f elrek-system_go_dev
                    fi
                    
                    if [ "\$(docker images -q sc4n1a471/elrek-system_go:$version-dev)" ]; then
                        docker rmi -f sc4n1a471/elrek-system_go:$version-dev
                    fi
                    exit
                    """

                    sh """
                    terraform init

                    terraform apply \
                        -var="container_version=$version-dev" \
                        -var="env=dev" \
                        -var="db_username=$DB_USERNAME" \
                        -var="db_name=$DB_NAME_DEV" \
                        -var="db_password=$DB_PASSWORD" \
                        -var="db_host=$DB_HOST" \
                        -var="db_port=$DB_PORT" \
                        -var="domain=$DOMAIN" \
                        -var="frontend_url=$FRONTEND_URL_DEV" \
                        -var="backend_url=$BACKEND_URL_DEV" \
                        -var="ssh_host=$SSH_HOST" \
                        -var="graylog_host=$GRAYLOG_HOST_DEV" \
                        -auto-approve
                    """
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
                    ssh -tt $SSH_HOST << EOF
                    if [ "\$(docker ps -a -q -f name=elrek-system_go_prod)" ]; then
                        docker rm -f elrek-system_go_prod
                    fi
                    
                    if [ "\$(docker images -q sc4n1a471/elrek-system_go:$version)" ]; then
                        docker rmi -f sc4n1a471/elrek-system_go:$version
                    fi
                    docker image rm sc4n1a471/elrek-system_go:$version
                    exit
                    """

                    sh """
                    terraform init

                    terraform apply \
                        -var="container_version=\$version" \
                        -var="env=prod" \
                        -var="db_username=$DB_USERNAME" \
                        -var="db_name=$DB_NAME_PROD" \
                        -var="db_password=$DB_PASSWORD" \
                        -var="db_host=$DB_HOST" \
                        -var="db_port=$DB_PORT" \
                        -var="domain=$DOMAIN" \
                        -var="frontend_url=$FRONTEND_URL_PROD" \
                        -var="backend_url=$BACKEND_URL_PROD" \
                        -var="ssh_host=$SSH_HOST" \
                        -var="graylog_host=$GRAYLOG_HOST_PROD" \
                        -auto-approve
                    """
                }
            }
        }
    }
    post {
        always {
            cleanWs()
            echo "Cleaning docker images"
            sh "docker rmi -f sc4n1a471/elrek-system_go"
            sh "docker image prune -f"
        }
    }
}