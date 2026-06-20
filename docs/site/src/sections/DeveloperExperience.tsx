export default function DeveloperExperience() {
  return (
    <section className="py-32">
      <div className="container mx-auto px-6 max-w-container">
        <div className="mb-16">
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-4">API First.</h2>
          <p className="text-text-muted text-lg max-w-2xl">
            Integrate Onyx into your internal developer platform using our strongly-typed gRPC and REST APIs.
          </p>
        </div>
        <div className="rounded-xl overflow-hidden border border-white/10 bg-black">
          <div className="flex border-b border-white/10 bg-white/5 px-4 text-sm text-text-muted">
            <div className="px-4 py-3 border-b-2 border-blue-500 text-white">main.go</div>
            <div className="px-4 py-3 border-b-2 border-transparent">deploy.py</div>
          </div>
          <div className="p-6 font-mono text-sm leading-relaxed overflow-x-auto text-blue-200">
<pre><code>{`import (
  "context"
  "github.com/onyx/sdk-go/onyx"
)

func main() {
  client := onyx.NewClient(onyx.WithAPIKey("sk_..."))

  job, err := client.Jobs.Create(context.Background(), &onyx.JobRequest{
    Name:    "llama-3-finetune",
    Image:   "pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime",
    GPUs:    64,
    Type:    onyx.GPUTypeH100,
    Command: []string{"python", "train.py"},
  })

  client.Jobs.WaitUntilRunning(context.Background(), job.ID)
}`}</code></pre>
          </div>
        </div>
      </div>
    </section>
  )
}
