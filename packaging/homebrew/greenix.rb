# Homebrew formula for the Greenix CLI.
# Publisher: Edgicode Limited
#
# Copy this file into a tap repository (github.com/greenixdb/homebrew-tap) as
# Formula/greenix.rb and replace the sha256 values with the ones published in
# the release's SHA256SUMS file. Users then install with:
#
#   brew install greenixdb/tap/greenix
#
# Homebrew downloads do not carry the macOS quarantine flag, so the binary runs
# without Gatekeeper prompts even before Apple notarization is in place.
class Greenix < Formula
  desc "Build, deploy and manage Greenix Studio projects"
  homepage "https://github.com/greenixdb/cli"
  version "0.1.0"
  license :cannot_represent # Proprietary - Edgicode Limited

  on_macos do
    url "https://github.com/greenixdb/cli/releases/download/v#{version}/greenix-macos-universal.tar.gz"
    sha256 "REPLACE_WITH_SHA256_FROM_SHA256SUMS"
  end

  on_linux do
    on_intel do
      url "https://github.com/greenixdb/cli/releases/download/v#{version}/greenix-linux-x64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_FROM_SHA256SUMS"
    end
    on_arm do
      url "https://github.com/greenixdb/cli/releases/download/v#{version}/greenix-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_FROM_SHA256SUMS"
    end
  end

  def install
    bin.install "greenix"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/greenix --version")
  end
end
