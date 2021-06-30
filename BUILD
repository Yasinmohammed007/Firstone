load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")
load("@bazel_sonarqube//:defs.bzl", "sonarqube")

# gazelle:prefix github.com/SophosNSG/paloma-config-service
# gazelle:proto disable_global
gazelle(name = "gazelle")

container_image(
    name = "config-service-image",
    base = "@alpine_amd64//image",
    entrypoint = ["/paloma-config-service"],
    files = ["//:paloma-config-service"],
    stamp = True,
    visibility = ["//visibility:public"],
)

container_push(
    name = "config-service-push",
    format = "Docker",
    image = ":config-service-image",
    registry = "$(docker_registry)",
    repository = "$(repository)",
    tag = "$(version)",
    visibility = ["//visibility:public"],
)

go_library(
    name = "paloma-config-service_lib",
    srcs = ["main.go"],
    importpath = "github.com/SophosNSG/paloma-config-service",
    visibility = ["//visibility:private"],
    deps = [
        "//src",
        "@com_github_gorilla_mux//:mux",
        "@com_github_opentracing_opentracing_go//:opentracing-go",
        "@com_github_opentracing_opentracing_go//ext",
        "@com_github_opentracing_opentracing_go//log",
        "@com_github_prnvkv_my_nats//pkg/util",
        "@com_github_sirupsen_logrus//:logrus",
        "@com_github_uber_jaeger_client_go//:jaeger-client-go",
        "@com_github_uber_jaeger_client_go//config",
    ],
)

go_binary(
    name = "paloma-config-service",
    embed = [":paloma-config-service_lib"],
    goarch = "amd64",  # remove this if you are building on linux
    goos = "linux",  # remove this if you are building on linux
    visibility = ["//visibility:public"],
)

go_test(
    name = "paloma-config-service_test",
    srcs = ["config_test.go"],
    embed = [":paloma-config-service_lib"],
    deps = ["//src"],
)
# filegroup(
#     name = "git",
#     srcs = glob(
#         [".git/**"],
#         exclude = [".git/**/*[*"],  # gitk creates temp files with []
#     ),
#     tags = ["manual"],
# )
# filegroup(
#     name = "coverage_report",
#     # srcs = ["bazel-out/darwin-fastbuild/testlogs/paloma-config-service_test/coverage.dat"],  # Created manually
#     srcs = ["coverage.dat"],  # Created manually
#     tags = ["manual"],
#     visibility = ["//visibility:public"],
# )

# filegroup(
#     name = "test_reports",
#     srcs = glob(["bazel-testlogs/**/test.xml"]), # Created manually
#     tags = ["manual"],
#     visibility = ["//visibility:public"],
# )

# sonarqube(
#     name = "sonarqube",
#     project_key = "sophos:nsg:central:paloma",
#     coverage_report = ":coverage_report",
#     project_name = "sample-app",
#     srcs = [
#         "main.go",
#     ],
#     targets = [

#     ],
#     scm_info = [":git"],
#     tags = ["manual"],
#     visibility = ["//visibility:public"],
# )
