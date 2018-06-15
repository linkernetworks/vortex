pipeline {
    agent {
        docker {
            image "linkernetworks/jenkins-docker-builder:ubuntu16.04"
            args "--privileged --group-add docker"
            alwaysPull true
        }
    }
    post {
        always {
            sh "make clean"
        }
    }
    stages {
        stage('Prepare') {
            steps {
                sh "go get -u github.com/kardianos/govendor"
                sh "make pre-build"
                sh "docker run -itd -p 27017:27017 --name mongo mongo"
            }
        }
        stage('Build') {
            steps {
                sh "make build"
            }
        }
        stage('Test') {
            steps {
                sh "make src.test-coverage"
            }
        }
    }
}