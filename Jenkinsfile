stage("build") {
  scm checkout
  docker.image("golang:1.7").inside {
    sh "go build"
  }
}
