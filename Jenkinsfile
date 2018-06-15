pipeline {
    agent none
    stages {
        stage('Prepare') {
            docker {
                image "linkernetworks/jenkins-docker-builder:ubuntu16.04"
                args "--privileged --group-add docker"
                alwaysPull true
            }
            steps {
                echo "hi jenkins"
            }
        }
    }
}