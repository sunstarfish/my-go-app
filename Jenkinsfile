pipeline {
    // agent any
     // 使用包含 Go 的 Docker 镜像作为代理
    // agent {
    //     docker {
    //         image 'golang:1.23-alpine'  // 官方 Go 镜像
    //         args '-v /var/run/docker.sock:/var/run/docker.sock'  // 允许在容器内使用 Docker
    //     }
    // }
    agent {
        docker {
            // 使用包含 Go 和 Docker 的镜像
            image 'golang:1.23-alpine'
            
            // 挂载 Docker socket 并安装 Docker CLI
            args '''
                -v /var/run/docker.sock:/var/run/docker.sock
                -v /usr/bin/docker:/usr/bin/docker
                -v /usr/libexec/docker:/usr/libexec/docker
            '''
            
            // 或者在容器内安装 Docker CLI
            args '''
                -v /var/run/docker.sock:/var/run/docker.sock
                && apk add --no-cache docker-cli
            '''
        }
    }

    environment {
        // 使用你自己的信息替换
        DOCKER_IMAGE = 'sunstarfish/my-go-app'  // 你的 Docker Hub 用户名/镜像名
        GIT_REPO = 'https://github.com/sunstarfish/my-go-app.git'
        
        // 可选：添加更多环境变量
        GO_VERSION = '1.23'
          // 设置国内 Go 代理（关键！）
        GOPROXY = 'https://goproxy.cn,direct'
        GOSUMDB = 'off'  // 关闭校验，国内网络可能无法访问 sum.golang.org
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
                sh '''
                    echo "=== Go 环境 ==="
                    go version
                    go env | grep -E "(GOPROXY|GOSUMDB|GOPATH|GO111MODULE)"
                    
                    echo "=== 网络测试 ==="
                    echo "测试网络连接..."
                    ping -c 2 goproxy.cn || echo "ping 测试失败"
                    curl -I https://goproxy.cn || echo "curl 测试失败"

                    echo "=== Docker 检查 ==="
                    # 尝试安装 Docker CLI
                    if ! command -v docker &> /dev/null; then
                        echo "安装 Docker CLI..."
                        apk add --no-cache docker-cli || echo "Docker CLI 安装失败"
                    fi
                    
                    docker --version || echo "Docker 不可用"
                    docker info 2>/dev/null || echo "无法连接 Docker daemon"
                    
                    echo "=== 当前目录 ==="
                    pwd
                    ls -la
                '''
            }
        }

        stage('Build and Test') {
            steps {
                sh '''
                    echo "=== 当前目录 ==="
                    pwd
                    ls -la
                    
                    echo "=== 检查 go.mod ==="
                    cat go.mod || echo "go.mod 不存在"
                    
                    echo "=== 下载依赖 ==="
                    go env GOPROXY
                    go mod download -x

                    echo "清理和下载依赖..."
                    go mod tidy
                    
                    echo "构建项目..."
                    go build -v ./...
                    
                    echo "运行测试..."
                    go test -v ./...

                    echo "=== 查看 vendor（如果有） ==="
                    ls -la vendor/ 2>/dev/null || echo "没有 vendor 目录"
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