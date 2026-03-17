class GoogleScholarMcp < Formula
  desc "Google Scholar MCP server written in Go"
  homepage "https://github.com/bingshuoguo/google-scholar-mcp"
  version "0.1.1"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/bingshuoguo/google-scholar-mcp/releases/download/v0.1.1/google-scholar-mcp_0.1.1_darwin_arm64.tar.gz"
      sha256 "27e3dd405405c041618098cc38ef90a0c0a6f7f9d1778106a76ea235c8bf0ead"
    else
      url "https://github.com/bingshuoguo/google-scholar-mcp/releases/download/v0.1.1/google-scholar-mcp_0.1.1_darwin_amd64.tar.gz"
      sha256 "9a12ec35a85cad43e474a133d3c9e6504ab34171ca6aba7cbd19985f70057872"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/bingshuoguo/google-scholar-mcp/releases/download/v0.1.1/google-scholar-mcp_0.1.1_linux_arm64.tar.gz"
      sha256 "f9887731259b5688e2c25f431362fac9eeb65924a5862eebc1da675a146353f3"
    else
      url "https://github.com/bingshuoguo/google-scholar-mcp/releases/download/v0.1.1/google-scholar-mcp_0.1.1_linux_amd64.tar.gz"
      sha256 "153e6931f0fcfde4bc862f18487a769567d36fe29fadbefdc48720fbd8415649"
    end
  end

  def install
    bin.install "google-scholar-mcp"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/google-scholar-mcp --version")
  end
end
