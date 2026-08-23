# Homebrew formula for Arcane Translator.
#
# This belongs in a tap repository (homebrew-tap) as
# Formula/arcane-translator.rb, so users can run:
#
#   brew install DebanganThakuria/tap/arcane-translator
#
# Before the first release: create a GitHub release tagged v1.0.0, then fill in
# `sha256` with the output of
#
#   curl -sL https://github.com/DebanganThakuria/arcane-translator/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
class ArcaneTranslator < Formula
  desc "Read Chinese, Korean and Japanese webnovels translated into English"
  homepage "https://github.com/DebanganThakuria/arcane-translator"
  url "https://github.com/DebanganThakuria/arcane-translator/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "985969c8a6a55c1a995c392aeb43100da1950380f8e3decbca0508a5fe8add63"
  license "Apache-2.0"
  head "https://github.com/DebanganThakuria/arcane-translator.git", branch: "master"

  depends_on "go" => :build
  depends_on "node" => :build

  def install
    # The frontend is built first and shipped as static files; the Go binary
    # serves them, so a running install is a single process on a single port.
    cd "web" do
      system "npm", "ci"
      system "npm", "run", "build"
    end

    cd "backend" do
      ldflags = "-s -w"
      system "go", "build", *std_go_args(output: bin/"arcane-translator", ldflags: ldflags)
    end

    (pkgshare/"web").install Dir["web/build/*"]
    pkgshare.install ".env.example"
  end

  def post_install
    (var/"arcane-translator").mkpath
  end

  service do
    run [opt_bin/"arcane-translator"]
    keep_alive true
    working_dir var/"arcane-translator"
    environment_variables ARCANE_WEB_DIR: "#{opt_pkgshare}/web",
                          ARCANE_DB_PATH: "#{var}/arcane-translator/data.db"
    log_path var/"log/arcane-translator.log"
    error_log_path var/"log/arcane-translator.log"
  end

  def caveats
    <<~EOS
      Arcane Translator needs a translation provider before it will start.

      Create #{Dir.home}/.arcane-translator/.env with at least:

        ARCANE_LLM_PROVIDER=claude
        AWS_REGION=us-east-1
        ARCANE_BEDROCK_MODEL_ID=<model id or inference profile ARN>

      AWS credentials are read from the standard AWS chain. For Gemini or an
      OpenAI-compatible endpoint see:
        #{opt_pkgshare}/.env.example

      Then run it in the foreground:
        arcane-translator

      Or in the background:
        brew services start arcane-translator

      The app is served at http://localhost:8088

      It has no authentication and binds to all interfaces so you can open it
      from a phone on the same network. Do not expose it to the internet.
    EOS
  end

  test do
    # Without a provider configured the server must refuse to start and say so,
    # rather than starting in a broken state.
    output = shell_output("#{bin}/arcane-translator 2>&1", 1)
    assert_match "Configuration error", output
  end
end
