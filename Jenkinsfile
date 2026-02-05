pipeline {
    // agent any
     // 使用包含 Go 的 Docker 镜像作为代理
    agent {
        docker {
            image 'golang:1.23-alpine'  // 官方 Go 镜像
            args '-v /var/run/docker.sock:/var/run/docker.sock'  // 允许在容器内使用 Docker
        }
    }

    environment {
        // 使用你自己的信息替换
        DOCKER_IMAGE = 'sunstarfish/my-go-app'  // 你的 Docker Hub 用户名/镜像名
        GIT_REPO = 'https://github.com/sunstarfish/my-go-app.git'
        
        // 可选：添加更多环境变量
        GO_VERSION = '1.23'
        DOCKER_REGISTRY = 'https://index.docker.io/v1/'
    }

    stages {
        stage('Checkout') {
            steps {
                echo "开始检出代码，仓库: ${GIT_REPO}"
                git url: "${GIT_REPO}", 
                     branch: 'master', 
                     credentialsId: 'GitHub'  // 确保 Jenkins 中有这个凭据
            }
        }

        stage('Setup') {
            steps {
                echo "当前构建ID: ${BUILD_ID}"
                echo "构建号: ${BUILD_NUMBER}"
                echo "工作空间: ${WORKSPACE}"
                
                // 检查 Go 版本
                sh 'go version'
            }
        }

        stage('Build and Test') {
            steps {
                sh '''
                    echo "清理和下载依赖..."
                    go mod tidy
                    
                    echo "构建项目..."
                    go build -v ./...
                    
                    echo "运行测试..."
                    go test -v ./...
                '''
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    echo "开始构建 Docker 镜像..."
                    echo "镜像标签: ${DOCKER_IMAGE}:${BUILD_ID}"
                    
                    // 构建 Docker 镜像
                    docker.build("${DOCKER_IMAGE}:${BUILD_ID}")
                    
                    echo "Docker 镜像构建完成"
                }
            }
        }

        stage('Push Docker Image') {
            steps {
                script {
                    echo "开始推送 Docker 镜像到 Docker Hub..."
                    
                    // 登录 Docker Hub 并推送
                    docker.withRegistry('https://index.docker.io/v1/', 'DockerHub') {
                        // 推送带构建ID的版本
                        docker.image("${DOCKER_IMAGE}:${BUILD_ID}").push()
                        echo "已推送镜像: ${DOCKER_IMAGE}:${BUILD_ID}"
                        
                        // 可选：推送 latest 标签
                        docker.image("${DOCKER_IMAGE}:${BUILD_ID}").push('latest')
                        echo "已推送 latest 标签"
                    }
                    
                    echo "镜像推送完成"
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    echo "开始部署..."
                    
                    // 检查是否有 docker-compose.yml 文件
                    sh 'ls -la docker-compose.yml || echo "未找到 docker-compose.yml"'
                    
                    // 停止并启动容器
                    sh '''
                        # 停止旧容器（如果存在）
                        docker compose down || true
                        
                        # 拉取最新镜像
                        docker pull ${DOCKER_IMAGE}:${BUILD_ID}
                        
                        # 启动新容器
                        docker compose up -d
                        
                        # 检查容器状态
                        sleep 5
                        docker ps
                    '''
                    
                    echo "部署完成"
                }
            }
        }
    }

    post {
        success {
            echo "流水线执行成功！"
            echo "镜像地址: ${DOCKER_IMAGE}:${BUILD_ID}"
        }
        failure {
            echo "流水线执行失败！"
        }
        always {
            echo "清理工作空间..."
            cleanWs()
            
            // 可选：清理 Docker 临时文件
            // sh 'docker system prune -f'
        }
    }
}