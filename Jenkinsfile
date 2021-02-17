#!groovy
@Library(['github.com/cloudogu/zalenium-build-lib@v2.1.0'])
import com.cloudogu.ces.zaleniumbuildlib.*

node('docker') {

    branch = "${env.BRANCH_NAME}"

    stage('Checkout') {
        checkout scm
    }

    def redmineImage = docker.image('redmine:4.1.1-alpine')
    def redmineContainerName = "${JOB_BASE_NAME}-${BUILD_NUMBER}".replaceAll("\\/|%2[fF]", "-")
    withDockerNetwork { buildnetwork ->
        redmineImage.withRun("--network ${buildnetwork} --name ${redmineContainerName} -p 8080:3000") {

            docker.image('golang:1.14.13').inside("--network ${buildnetwork} -e HOME=/tmp") {

                stage('Build') {
                    make 'clean package checksum'
                    archiveArtifacts 'target/*'
                }

                stage('Unit Test') {
                    make 'unit-test'
                    junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
                }

                stage('Static Analysis') {
                    make 'static-analysis'
                }

                stage('Acceptance Tests') {
                    sh "REDMINE_URL=http://${redmineContainerName}:8080/ make testacc"
                    archiveArtifacts 'target/acceptance-tests/*.out'
                }

            }
        }
    }
    stage('SonarQube') {
        def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
        withSonarQubeEnv {
            if (branch == "main") {
                echo "This branch has been detected as the main branch."
                sh "${scannerHome}/bin/sonar-scanner"
            } else if (branch == "develop") {
                echo "This branch has been detected as the develop branch."
                sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=master"
            } else if (env.CHANGE_TARGET) {
                echo "This branch has been detected as a pull request."
                sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.CHANGE_BRANCH}-PR${env.CHANGE_ID} -Dsonar.branch.target=${env.CHANGE_TARGET}"
            } else if (branch.startsWith("feature/")) {
                echo "This branch has been detected as a feature branch."
                sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
            }
        }
        timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
            def qGate = waitForQualityGate()
            if (qGate.status != 'OK') {
                unstable("Pipeline unstable due to SonarQube quality gate failure")
            }
        }
    }
}

void make(String goal) {
    sh "make ${goal}"
}