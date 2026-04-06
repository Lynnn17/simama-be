pipeline {
    agent any

    environment {
        DOCKER_IMAGE_NAME = '/lms/lms-be'
        PATH_REPOSITORY = "/lms/lms-be.git"
        APPNAME = "lms-be"
        TIMESTAMP = "${currentBuild.getTimeInMillis()}"
        IMAGE_TAG = "${TIMESTAMP}-${BUILD_ID}"
        FULL_IMAGE_NAME = "${REGISTRY_URL}${DOCKER_IMAGE_NAME}:${IMAGE_TAG}"
        GIT_EMAIL = "${GITLAB_USER_EMAIL}"
        GIT_NAME = "${GITLAB_USER_NAME}"
    }

     stages {
        stage('SonarQube Analysis') {
            steps {
                script {
                    def scannerHome = tool 'SonarQube' 
                    withSonarQubeEnv('sonar-server') {
                        sh """
                            ${scannerHome}/bin/sonar-scanner \
                            -Dproject.settings=sonar-project.properties
                        """
                    }
                }
            }
        }

        stage('Build Image Docker') {
            steps {
                script {
                    sh """
                        docker build -t ${FULL_IMAGE_NAME} .
                        docker push ${FULL_IMAGE_NAME}
                    """
                }
            }
        }

        stage('Update Image Tag in YAML') {
            steps {
                script {
                    // Update image tag in deployment and kustomization YAML files
                    sh """
                        sed -i 's|image:.*|image: ${FULL_IMAGE_NAME}|g' ${PATH_KUSTOMIZE}/base/${APPNAME}/deployment.yaml
                        sed -i 's|newName:.*|newName: ${REGISTRY_URL}${DOCKER_IMAGE_NAME}|g' ${PATH_KUSTOMIZE}/overlays/production/${APPNAME}/kustomization.yaml
                        sed -i 's|newTag:.*|newTag: ${IMAGE_TAG}|g' ${PATH_KUSTOMIZE}/overlays/production/${APPNAME}/kustomization.yaml
                    """

                    // Print the updated files for verification
                    sh """
                        cat ${PATH_KUSTOMIZE}/base/${APPNAME}/deployment.yaml
                        cat ${PATH_KUSTOMIZE}/overlays/production/${APPNAME}/kustomization.yaml
                    """
                }
            }
        }

         stage('Commit and Push Changes') {
            steps {
                script {
                    // Navigate to the kustomize directory and push changes to Git
                    dir("${PATH_KUSTOMIZE}") {
                        sh """
                            git config --global user.email ${GIT_EMAIL}
                            git config --global user.name ${GIT_NAME}
                            git pull origin main
                            git add .
                            git commit -m "Update Docker Image to ${IMAGE_TAG}"
                            git push origin main
                        """
                    }
                }
            }
        }
    }

   post {
        always {
            script {
                def status = currentBuild.currentResult
                def message = (status == 'SUCCESS') 
                    ? "Pipeline '${currentBuild.fullDisplayName}' MasyaAllah, Alhamdulillah Build Sukses :white_check_mark:" 
                    : "Pipeline '${currentBuild.fullDisplayName}' Innalillahi, Build Gagal  :x:"

                discordSend title: "Build And Push Docker Job $JOB_NAME",
                            link: "${BUILD_URL}console",
                            webhookURL: "$DISCORD_WEBHOOK",
                            result: status, 
                            description: message,
                            footer: "${APPNAME}"
                
               // Notifikasi untuk bug dari SonarQube
                def sonarIssues = sh(script: "curl -s -X GET 'http://10.10.10.11:9000/api/issues/search?componentKeys=project-gsi&severities=BLOCKER,CRITICAL'", returnStdout: true)
                
                // Periksa apakah sonarIssues tidak kosong dan merupakan JSON yang valid
                if (sonarIssues) {
                    def issuesData
                    try {
                        issuesData = readJSON(text: sonarIssues)
                    } catch (Exception e) {
                        echo "Gagal untuk mengurai respons SonarQube: ${e.message}"
                        return // Keluar jika penguraian gagal
                    }
                    
                    def issuesCount = issuesData.total ?: 0 // Akses total dengan aman

                    if (issuesCount > 0) {
                        def issueMessage = "SonarQube menemukan **${issuesCount}** masalah di ${APPNAME}. Silakan periksa SonarQube untuk detail."
                        sh """
                            curl -H "Content-Type: application/json" -X POST -d '{
                                "content": "${issueMessage}"
                            }' ${DISCORD_WEBHOOK}
                        """
                    }
                } else {
                    echo "Tidak ada masalah ditemukan dari SonarQube atau responsnya kosong."
                }
            }
        }
    }
}