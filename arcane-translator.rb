class ArcaneTranslator < Formula
  desc "A web application for translating novels with AI"
  homepage "https://github.com/DebanganDaemon/arcane-translator"
  url "https://github.com/DebanganDaemon/arcane-translator/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "" # You'll need to calculate this after creating the release
  license "MIT"

  depends_on "go" => :build
  depends_on "node" => :build

  def install
    ENV.deparallelize

    # Install frontend dependencies and build
    system "npm", "install"
    system "npm", "run", "build"

    # Install backend dependencies and build
    cd "backend" do
      system "go", "mod", "download"
      system "go", "build", "-o", "arcane-translator-backend"
    end

    # Install the binary
    bin.install "arcane-translator"
    bin.install "backend/arcane-translator-backend"
    
    # Install the built frontend
    pkgshare.install "dist"
  end

  def post_install
    # Create a plist file for launchd
    (prefix/"homebrew.mxcl.arcane-translator.plist").write plist
    (prefix/"homebrew.mxcl.arcane-translator.plist").chmod 0644
  end

  def caveats
    <<~EOS
      To start Arcane Translator, run:
        arcane-translator

      The application will be available at http://localhost:8088
    EOS
  end

  test do
    system "#{bin}/arcane-translator-backend", "--version"
  end
end
