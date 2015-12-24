# terraform-provisioner-waffles

A [Waffles](http://waffles.terrarum.net) provisioner for Terraform

## Usage

```hcl
provisioner "waffles" {
  host = "${openstack_compute_instance_v2.scratch.access_ip_v6}"
  site_directory = "~/waffles"
  role = "scratch"
}
```

## Attributes

* `debug`: Run Waffles in debug mode. Optional.
* `host`: The address of the remote host. Required.
* `private_key`: The path of the private key to connect to the remote host. Optional. Defaults to `~/.ssh/id_rsa`.
* `remote_dir`: The path on the remote host to copy the Waffles profiles to. Optional. Defaults to `/etc/waffles`.
* `retry`: The number of SSH retry attempts. Optional. Defaults to 5.
* `role`: The role to apply to the remote host. Required.
* `site_directory`: The path to the local `WAFFLES_SITE_DIR`. Required.
* `sudo`: Whether or not to run Waffles remotely using `sudo`. Optional. Defaults to false.
* `user`: The user to connect to the remote host. Optional. Defaults to `root`.
* `waffles_exec`: The path to `waffles.sh`. Optional. Defaults to `/etc/waffles/waffles.sh`.
* `wait`: The amount of time in seconds to wait between SSH attempts. Optional Defaults to 5 seconds.

## Installation

1. Grab the latest release from the [releases](https://github.com/jtopjian/terraform-provisioner-waffles/releases) page.
2. Copy the binary to the same location as the other Terraform executables.

## Building

```shell
$ go get github.com/jtopjian/terraform-provisioner-waffles
$ cd $GOPATH/src/github.com/jtopjian/terraform-provisioner-waffles
$ go build -v -o ~/path/to/terraform/terraform-provisioner-waffles .
```
