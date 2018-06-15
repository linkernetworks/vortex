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
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        sh "make pre-build"
                    }
                }
            }
        }
        stage('Build') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        sh "make build"
                    }
                }
            }
        }
        stage('Test') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/vortex") {
                        docker.image('mongo').withRun('-p 27017:27017') { c ->
                            sh "make src.test-coverage"
                        }
                    }
                }
            }
        }
    }
}