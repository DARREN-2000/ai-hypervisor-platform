import * as Accordion from '@radix-ui/react-accordion'
import { ChevronDown } from 'lucide-react'

export default function FAQ() {
  const faqs = [
    { q: "Do I need to migrate off Kubernetes?", a: "No. Onyx runs alongside or on top of your existing Kubernetes clusters. It acts as a specialized scheduler for ML workloads." },
    { q: "How is billing handled for multi-cloud?", a: "Onyx provides a unified billing dashboard. You pay your cloud providers directly, and Onyx charges a small orchestration fee per GPU hour." },
    { q: "Is Onyx open source?", a: "The core scheduling and eBPF networking engine is open source. The enterprise control plane and RBAC features are commercial." }
  ]

  return (
    <section className="py-32">
      <div className="container mx-auto px-6 max-w-container max-w-3xl">
        <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-12 text-center">Frequently asked questions</h2>
        <Accordion.Root type="single" collapsible className="space-y-4">
          {faqs.map((faq, i) => (
            <Accordion.Item key={i} value={`item-${i}`} className="border-b border-white/10 overflow-hidden">
              <Accordion.Header>
                <Accordion.Trigger className="flex w-full items-center justify-between py-6 text-left text-lg font-medium transition-all hover:text-white text-text-muted [&[data-state=open]>svg]:rotate-180 [&[data-state=open]]:text-white">
                  {faq.q}
                  <ChevronDown className="h-5 w-5 shrink-0 transition-transform duration-200" />
                </Accordion.Trigger>
              </Accordion.Header>
              <Accordion.Content className="overflow-hidden text-sm text-text-muted data-[state=closed]:animate-accordion-up data-[state=open]:animate-accordion-down mb-6">
                {faq.a}
              </Accordion.Content>
            </Accordion.Item>
          ))}
        </Accordion.Root>
      </div>
    </section>
  )
}
