pipeline {
    agent any

    environment {
        DOCKER_IMAGE = 'sunstarfish/my-go-app'
        GIT_REPO = 'https://github.com/sunstarfish/my-go-app.git'
    }

    stages {
        stage('Checkout') {
            steps {
                git url: "${GIT_REPO}", branch: 'master', credentialsId: 'GitHub'
            }
        }

        stage('Build and Test') {
            steps {
                sh 'go mod tidy'
                sh 'go build ./...'
                sh 'go test ./...'
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    docker.build("${DOCKER_IMAGE}:${env.BUILD_ID}")
                }
            }
        }

        stage('Push Docker Image') {
            steps {
                script {
                    docker.withRegistry('https://index.docker.io/v1/', 'DockerHub') {
                        docker.image("${DOCKER_IMAGE}:${env.BUILD_ID}").push()
                        docker.image("${DOCKER_IMAGE}:${env.BUILD_ID}").push('latest')
                    }
                }
            }
        }

        stage('Deploy') {
            steps {
                sh 'docker compose down || true'  // 停止旧容器
                sh 'docker compose up -d'  // 启动新版本
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}