class Ghb < Formula
  desc "A tool for performing backups and deletions of Github gists and starred repositories"
  homepage ""
  url "https://github.com/mmercedes/ghb/releases/download/v0.1.1/ghb_0.1.1_darwin_amd64.tar.gz"
  version "0.1.1"
  sha256 "44d48173abc55bf577dc9a6a710a2387486ac494c562d8548da24106448c15bd"

  def install
    bin.install "ghb"
  end
end
