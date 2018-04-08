class Ghb < Formula
  desc "A tool for performing backups and deletions of Github gists and starred repositories"
  homepage ""
  url "https://github.com/mmercedes/ghb/releases/download/v0.1.0/ghb_0.1.0_darwin_amd64.tar.gz"
  version "0.1.0"
  sha256 "99d70bafb74acd979a73a3b0ebe075d257ba92b01f2ad55128d86d59ffad506d"

  def install
    bin.install "ghb"
  end
end
