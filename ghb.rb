class Ghb < Formula
  desc "A tool for performing backups and deletions of Github gists and starred repositories"
  homepage ""
  url "https://github.com/mmercedes/ghb/releases/download/v0.1.0/ghb_0.1.0_darwin_amd64.tar.gz"
  version "0.1.0"
  sha256 "999d6994b8efaa1d2abd166821e67e6943694b5d91b9cf0663256b360d42a9b8"

  def install
    bin.install "ghb"
  end
end
