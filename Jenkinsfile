pipeline {
    agent {
        dockerfile {
            dir "src/github.com/linkernetworks/vortex/jenkins"
            args "--privileged --group-add docker"
        }
    }
    post {
        always {
            dir ("src/github.com/linkernetworks/vortex") {
                sh "make clean"
            }
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
                        sh "make src.test-coverage 2>&1 | tee >(go-junit-report > report.xml)"
                        junit "report.xml"
                        sh 'gocover-cobertura < build/src/coverage.txt > cobertura.xml'
                        cobertura coberturaReportFile: "cobertura.xml", failNoReports: true, failUnstable: true
                        publishHTML (target: [
                            allowMissing: true,
                            alwaysLinkToLastBuild: true,
                            keepAll: true,
                            reportDir: 'build/src',
                            reportFiles: 'coverage.html',
                            reportName: "GO cover report",
                            reportTitles: "GO cover report",
                            includes: "coverage.html"
                        ])
                    }
                }
            }
        }
    }
}