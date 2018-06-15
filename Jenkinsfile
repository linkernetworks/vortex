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
                script {
                    dir ("src/github.com/linkernetworks/vortex") {
                        withEnv([
                            "PATH+GO=${env.WORKSPACE}aaaa/bin",
                            "GOPATH=${env.WORKSPACE}",
                        ]) {

                            sh "go get -u github.com/kardianos/govendor"
                            sh "ls"
                            sh "pwd"
                            sh 'echo GOPATH=$GOPATH'
                            sh 'echo PATH=$PATH'
                            sh 'ls $GOPATH/bin'
                            sh "make pre-build"
                            sh "docker run -itd -p 27017:27017 --name mongo mongo"
                        }
                    }
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