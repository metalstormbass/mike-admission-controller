package main

deny[msg] {
    input.kind == "Pod"
    container := input.spec.containers[_]
    not endswith(container.image, ":latest")
    msg = sprintf("Container '%v' in Pod '%v' is not using 'latest' tag", [container.name, input.metadata.name])
}