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
                dir ("src/github.com/linkernetworks/vortex") {
                    sh "make pre-build"
                }
            }
        }
        stage('Build') {
            steps {
                withEnv([
                    "GOPATH=${env.WORKSPACE}",
                    "PATH+GO=${env.WORKSPACE}/bin",
                ]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        sh "make build"
                    }
                }
            }
        }
        stage('Test') {
            steps {
                withEnv([
                    "GOPATH=${env.WORKSPACE}",
                    "PATH+GO=${env.WORKSPACE}/bin",
                ]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        sh "make src.test-coverage"
                    }
                }
            }
        }
    }
}