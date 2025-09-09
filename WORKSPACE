workspace(name = "faceit_cli")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Download Go rules
http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://github.com/bazelbuild/rules_go/releases/download/v0.50.1/rules_go-v0.50.1.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

# Go rules setup
go_rules_dependencies()

go_register_toolchains(version = "1.23.0")

# Download Gazelle
http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.43.0/bazel-gazelle-v0.43.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_repository(
    name = "com_github_antihax_optional",
    importpath = "github.com/antihax/optional",
    sum = "h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_aymanbagabas_go_osc52_v2",
    importpath = "github.com/aymanbagabas/go-osc52/v2",
    sum = "h1:HwpRHbFMcZLEVr42D4p7XBqjyuxQH5SMiErDT4WkJ2k=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_aymanbagabas_go_udiff",
    importpath = "github.com/aymanbagabas/go-udiff",
    sum = "h1:TK0fH4MteXUDspT88n8CKzvK0X9O2xu9yQjWpi6yML8=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_bits_and_blooms_bitset",
    importpath = "github.com/bits-and-blooms/bitset",
    sum = "h1:Tquv9S8+SGaS3EhyA+up3FXzmkhxPGjQQCkcs2uw7w4=",
    version = "v1.22.0",
)

go_repository(
    name = "com_github_charmbracelet_bubbletea",
    importpath = "github.com/charmbracelet/bubbletea",
    sum = "h1:VkHIxPJQeDt0aFJIsVxw8BQdh/F/L2KKZGsK6et5taU=",
    version = "v1.3.6",
)

go_repository(
    name = "com_github_charmbracelet_colorprofile",
    importpath = "github.com/charmbracelet/colorprofile",
    sum = "h1:4pZI35227imm7yK2bGPcfpFEmuY1gc2YSTShr4iJBfs=",
    version = "v0.2.3-0.20250311203215-f60798e515dc",
)

go_repository(
    name = "com_github_charmbracelet_lipgloss",
    importpath = "github.com/charmbracelet/lipgloss",
    sum = "h1:vYXsiLHVkK7fp74RkV7b2kq9+zDLoEU4MZoFqR/noCY=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_charmbracelet_x_ansi",
    importpath = "github.com/charmbracelet/x/ansi",
    sum = "h1:BXt5DHS/MKF+LjuK4huWrC6NCvHtexww7dMayh6GXd0=",
    version = "v0.9.3",
)

go_repository(
    name = "com_github_charmbracelet_x_cellbuf",
    importpath = "github.com/charmbracelet/x/cellbuf",
    sum = "h1:vy0GVL4jeHEwG5YOXDmi86oYw2yuYUGqz6a8sLwg0X8=",
    version = "v0.0.13-0.20250311204145-2c3ea96c31dd",
)

go_repository(
    name = "com_github_charmbracelet_x_exp_golden",
    importpath = "github.com/charmbracelet/x/exp/golden",
    sum = "h1:G99klV19u0QnhiizODirwVksQB91TJKV/UaTnACcG30=",
    version = "v0.0.0-20240806155701-69247e0abc2a",
)

go_repository(
    name = "com_github_charmbracelet_x_term",
    importpath = "github.com/charmbracelet/x/term",
    sum = "h1:AQeHeLZ1OqSXhrAWpYUtZyX1T3zVxfpZuEQMIQaGIAQ=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_erikgeiser_coninput",
    importpath = "github.com/erikgeiser/coninput",
    sum = "h1:Y/CXytFA4m6baUTXGLOoWe4PQhGxaX0KpnayAqC48p4=",
    version = "v0.0.0-20211004153227-1c3628e74d0f",
)

go_repository(
    name = "com_github_joho_godotenv",
    importpath = "github.com/joho/godotenv",
    sum = "h1:7eLL/+HRGLY0ldzfGMeQkb7vMd0as4CfYvUVzLqw0N0=",
    version = "v1.5.1",
)

go_repository(
    name = "com_github_klauspost_compress",
    importpath = "github.com/klauspost/compress",
    sum = "h1:wKRjX6JRtDdrE9qwa4b/Cip7ACOshUI4smpCQanqjSY=",
    version = "v1.15.9",
)

go_repository(
    name = "com_github_lucasb_eyer_go_colorful",
    importpath = "github.com/lucasb-eyer/go-colorful",
    sum = "h1:1nnpGOrhyZZuNyfu1QjKiUICQ74+3FNCN69Aj6K7nkY=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_mattn_go_isatty",
    importpath = "github.com/mattn/go-isatty",
    sum = "h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=",
    version = "v0.0.20",
)

go_repository(
    name = "com_github_mattn_go_localereader",
    importpath = "github.com/mattn/go-localereader",
    sum = "h1:ygSAOl7ZXTx4RdPYinUpg6W99U8jWvWi9Ye2JC/oIi4=",
    version = "v0.0.1",
)

go_repository(
    name = "com_github_mattn_go_runewidth",
    importpath = "github.com/mattn/go-runewidth",
    sum = "h1:E5ScNMtiwvlvB5paMFdw9p4kSQzbXFikJ5SQO6TULQc=",
    version = "v0.0.16",
)

go_repository(
    name = "com_github_mconnat_go_faceit",
    importpath = "github.com/mconnat/go-faceit",
    sum = "h1:9haGsL+jZux0K4tx0TXNJZ9oXy/tKon31G+nsTIFZ6o=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_muesli_ansi",
    importpath = "github.com/muesli/ansi",
    sum = "h1:ZK8zHtRHOkbHy6Mmr5D264iyp3TiX5OmNcI5cIARiQI=",
    version = "v0.0.0-20230316100256-276c6243b2f6",
)

go_repository(
    name = "com_github_muesli_cancelreader",
    importpath = "github.com/muesli/cancelreader",
    sum = "h1:3I4Kt4BQjOR54NavqnDogx/MIoWBFa0StPA8ELUXHmA=",
    version = "v0.2.2",
)

go_repository(
    name = "com_github_muesli_termenv",
    importpath = "github.com/muesli/termenv",
    sum = "h1:S5AlUN9dENB57rsbnkPyfdGuWIlkmzJjbFf0Tf5FWUc=",
    version = "v0.16.0",
)

go_repository(
    name = "com_github_pierrec_lz4_v4",
    importpath = "github.com/pierrec/lz4/v4",
    sum = "h1:MO0/ucJhngq7299dKLwIMtgTfbkoSPF6AoMYDd8Q4q0=",
    version = "v4.1.15",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_rivo_uniseg",
    importpath = "github.com/rivo/uniseg",
    sum = "h1:WUdvkW8uEhrYfLC4ZzdpI2ztxP1I582+49Oc5Mq64VQ=",
    version = "v0.4.7",
)

go_repository(
    name = "com_github_segmentio_kafka_go",
    importpath = "github.com/segmentio/kafka-go",
    sum = "h1:GJiNX1d/g+kG6ljyJEoi9++PUMdXGAxb7JGPiDCuNmk=",
    version = "v0.4.49",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=",
    version = "v1.9.3",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:pSgiaMZlXftHpm5L7V1+rVB+AZJydKsMxsQBIJw4PKk=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_xdg_go_pbkdf2",
    importpath = "github.com/xdg-go/pbkdf2",
    sum = "h1:Su7DPu48wXMwC3bs7MCNG+z4FhcyEuz5dlvchbq0B0c=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_xdg_go_scram",
    importpath = "github.com/xdg-go/scram",
    sum = "h1:FHX5I5B4i4hKRVRBCFRxq1iQRej7WO3hhBuJf+UUySY=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_xdg_go_stringprep",
    importpath = "github.com/xdg-go/stringprep",
    sum = "h1:XLI/Ng3O1Atzq0oBs3TWm+5ZVgkq2aqdlvP9JtoZ6c8=",
    version = "v1.0.4",
)

go_repository(
    name = "com_github_xo_terminfo",
    importpath = "github.com/xo/terminfo",
    sum = "h1:JVG44RsyaB9T2KIHavMF/ppJZNG9ZpyihvCd0w101no=",
    version = "v0.0.0-20220910002029-abceb7e1c41e",
)

go_repository(
    name = "com_google_cloud_go_compute_metadata",
    importpath = "cloud.google.com/go/compute/metadata",
    sum = "h1:Tz+eQXMEqDIKRsmY3cHTL6FVaynIjX2QxYC4trgAKZc=",
    version = "v0.3.0",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
    version = "v0.0.0-20161208181325-20d25e280405",
)

go_repository(
    name = "in_gopkg_yaml_v3",
    importpath = "gopkg.in/yaml.v3",
    sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
    version = "v3.0.1",
)

go_repository(
    name = "org_golang_x_exp",
    importpath = "golang.org/x/exp",
    sum = "h1:MDc5xs78ZrZr3HMQugiXOAkSZtfTpbJLDr/lwfgO53E=",
    version = "v0.0.0-20220909182711-5c715a9e8561",
)

go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:zY54UmvipHiNd+pm+m0x9KhZ9hl1/7QNMyxXbc6ICqA=",
    version = "v0.17.0",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:vRMAPTMaeGqVhG5QyLJHqNDwecKTomGeqbnfZyKlBI8=",
    version = "v0.38.0",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    sum = "h1:dnDm7JmhM45NNpd8FDDeLhK6FwqbOf4MLCM9zb1BOHI=",
    version = "v0.30.0",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:KWH3jNZsfyT6xfAfKiz6MRNmd46ByHDYaZ7KSkCtdW8=",
    version = "v0.15.0",
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:q3i8TbbEz+JRD9ywIRlyRAQbM0qF7hu24q3teo2hbuw=",
    version = "v0.33.0",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:D71I7dUrlY+VX0gQShAThNGHFxZ13dGLBHQLVl1mJlY=",
    version = "v0.23.0",
)

go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    sum = "h1:vU5i/LfpvrRCpgM/VPfJLg5KjxD3E+hfT1SH+d9zLwg=",
    version = "v0.21.1-0.20240508182429-e35e4ccd0d2d",
)

# Gazelle setup
gazelle_dependencies()

# Configure Go proxy
load("@bazel_gazelle//:deps.bzl", "go_repository")
