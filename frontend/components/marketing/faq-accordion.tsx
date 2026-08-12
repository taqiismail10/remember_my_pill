import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

const faqs = [
  {
    question: "What is Remember My Pill?",
    answer:
      "Remember My Pill is a privacy-first medication support product in development. It's designed to turn prescription instructions into a clear routine: snap a prescription, understand it in plain language, and act on grouped reminders.",
  },
  {
    question: "Who is it for?",
    answer:
      "Anyone managing one or more daily medicines, and caregivers who help coordinate medication for family members.",
  },
  {
    question: "Is the mobile app available yet?",
    answer:
      "Not yet. This website is a pre-launch preview. The screens shown here illustrate the planned product design — they are not connected to real prescriptions or health data.",
  },
  {
    question: "What happens when I join the waitlist?",
    answer:
      "We store your email so we can notify you when early access opens. Joining does not create an account in the mobile product, and this preview does not process prescriptions or collect medication or health information.",
  },
  {
    question: "How does Remember My Pill approach privacy?",
    answer:
      "We collect only what the waitlist needs: your email. We don't use behavioral advertising trackers, and we don't ask for medication, diagnosis, or insurance details at this stage.",
  },
  {
    question: "Will I get a referral link or waitlist rank?",
    answer:
      "Not yet. Today's waitlist keeps things simple with just your email — referral links and rank are part of our product plans but are not part of this preview.",
  },
];

export function FaqAccordion() {
  return (
    <Accordion type="single" collapsible className="w-full">
      {faqs.map((faq, index) => (
        <AccordionItem key={faq.question} value={`faq-${index}`}>
          <AccordionTrigger>{faq.question}</AccordionTrigger>
          <AccordionContent>{faq.answer}</AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  );
}
