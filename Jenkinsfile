pipeline {
  agent none
  stages {
    stage('test') {
      agent {
        docker {
          image 'golang'
          label 'Slave-01'
        }

      }
      steps {
        sh 'GOFLAGS=-mod=vendor go test ./...'
      }
    }
    stage('build & push') {
      agent {
        label 'Slave-01'
      }
      steps {
        sh '''IMAGE_VERSION=$(date +%y%m%d%H%M)+B${BUILD_NUMBER}
IMAGE_TAG=192.168.5.82/grpc/grpc-demo:v0.0.${IMAGE_VERSION}
docker build -t ${IMAGE_TAG} -f ${HOME}/build/build-image/Dockerfile .
docker push IMAGE_TAG'''
      }
    }
  }
  environment {
    HOME = "${env.WORKSPACE}"
  }
}