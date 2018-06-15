pipeline {
    agent {
        dockerfile {
            dir "src/github.com/linkernetworks/vortex/jenkins"
            args "--privileged --group-add docker"
        }
    }
    post {
        always {
            sh "make clean"
        }
    }
    options {
        timestamps()
        timeout(time: 1, unit: 'HOURS')
        checkoutToSubdirectory('src/github.com/linkernetworks/vortex')
    }
    stages {
        stage('Prepare') {
            steps {
                withEnv([
                    "GOPATH=${env.WORKSPACE}",
                    "PATH+GO=${env.WORKSPACE}/bin",
                ]) {
                    sh "go get -u github.com/kardianos/govendor"
                    sh "ls"
                    sh "pwd"
                    sh "make pre-build"
                    sh "docker run -itd -p 27017:27017 --name mongo mongo"
                }
            }
        }
        stage('Build') {
            steps {
                withEnv([
                    "GOPATH=${env.WORKSPACE}",
                    "PATH+GO=${env.WORKSPACE}/bin",
                ]) {
                    sh "make build"
                }
            }
        }
        stage('Test') {
            steps {
                withEnv([
                    "GOPATH=${env.WORKSPACE}",
                    "PATH+GO=${env.WORKSPACE}/bin",
                ]) {
                    sh "make src.test-coverage"
                }
            }
        }
    }
}