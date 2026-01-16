import {
  Navbar,
  HeroSection,
  ProblemSection,
  SolutionSection,
  FeaturesSection,
  HowItWorksSection,
  UserRolesSection,
  StatisticsSection,
  CTASection,
  Footer,
} from "@/components/landing";

export default function LandingPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <main className="flex-1">
        <HeroSection />
        <ProblemSection />
        <SolutionSection />
        <FeaturesSection />
        <HowItWorksSection />
        <UserRolesSection />
        <StatisticsSection />
        <CTASection />
      </main>
      <Footer />
    </div>
  );
}
