class Ghb < Formula
  desc "A tool for performing backups and deletions of Github gists and starred repositories"
  homepage ""
  url "https://github.com/mmercedes/ghb/releases/download/v0.1.0/ghb_0.1.0_darwin_amd64.tar.gz"
  version "0.1.0"
  sha256 "f4b08d74c1782a45e2d1a57c921ec6c4756ed71e8e572a7cf3cbc7d77b975667"

  def install
    bin.install "ghb"
  end
end
