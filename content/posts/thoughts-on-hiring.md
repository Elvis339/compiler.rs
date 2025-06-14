+++
date = '2025-06-14T16:03:09+04:00'
draft = true
title = 'Thought on hiring'
description = 'How cognitive biases and low-validity environments sabotage hiring decisions and what we can do about it.'
+++

> Disclaimer: I'm not a trained psychologist, recruiter, or statistician. This article reflects my reasoning drawing mostly from Daniel Kahneman's "Thinking, Fast and Slow" and personal experience.

Most companies are terrible at hiring. Despite decades of research on human judgment and decision-making, most organizations continue to rely on processes that are barely better than random chance at predicting job performance.

I've experienced this firsthand through two contrasting but equally flawed approaches at well-known tech companies.

# The Current State of Software Engineering Hiring

The maximum number of interview rounds I've experienced was 9, and getting rejected after passing all of those rounds is a deeply disappointing and frustrating experience.

The usual hiring process follows this pattern:
- Initial call: if you're applying, recruiters assess you against dimensions listed in the job description. If you're being recruited, they're essentially doing the same evaluation in reverse.
- Technical Challenges: This stage typically includes LeetCode-style algorithmic questions, bug finding exercises, and integration rounds. While these are attempts to gauge technical competency by specific dimensions (as we'll discuss later in this article), the execution is fundamentally flawed. These challenges rarely reflect actual job responsibilities and often measure test-taking ability rather than engineering capability.
- System Design questions rarely have single "correct" solutions. Skilled interviewers understand this. They're not seeking the perfect architecture but evaluating whether you can manage time effectively, present working solutions, guide conversations toward meaningful technical discussions, and engage thoughtfully about tradeoffs. Yet many interviewers miss this entirely, instead looking for candidates to hit specific talking points or follow predetermined scripts.
- Behavioral Interviews: These sessions attempt to assess soft skills, cultural fit, and past performance through situational questions. However, they're almost entirely subjective, relying on the interviewer's ability to interpret stories and predict future behavior based on past examples. The evaluation criteria are rarely standardized, making these rounds particularly susceptible to bias and inconsistent scoring.

Each stage appears objective on the surface. After all, technical questions have right and wrong answers. But in practice, every stage involves substantial subjective judgment about communication style, thought processes, and cultural alignment. The frustrating reality is that you can excel at the supposedly objective technical components and still be rejected based on someone's gut feeling about whether you're "the right fit."

## The Psychology of Broken Hiring

The core issue in hiring doesn't seem to be technical rather psychological. When interviewers evaluate candidates, they unknowingly engage in what Kahneman calls "substitution"; replacing the difficult question "Will this person perform well in this role?" with easier questions like "Do they remind me of successful people I know?" or "Do they follow my predefined script for how this question should be answered?"

This substitution combines with the representativeness heuristic, where interviewers judge candidates based on how closely they match mental prototypes of "good employees." The problem is that interviews are artificial situations bearing little resemblance to actual job performance, making predictions fundamentally unreliable.

Perhaps most dangerously, subjective confidence in hiring decisions isn't a reasoned evaluation of accuracy it's merely a feeling reflecting the coherence of available information and the cognitive ease of processing it.

### Hiring is a low validity environment

Kahneman distinguishes between high-validity and low-validity environments based on their feedback mechanisms and pattern reliability.

High-validity environments feature regular patterns, immediate accurate feedback, and sufficient opportunity to learn from mistakes.

Low-validity environments lack these essential elements. Patterns are weak or nonexistent, feedback is delayed or unreliable, and signal-to-noise ratios are poor.

Hiring falls under the low-validity environment. Consider the feedback loop: you interview someone, make a hiring decision, and then... what? If you hire them, performance signals arrive months or years later, contaminated by countless variables such as team dynamics, management quality, project assignments, market conditions, onboarding effectiveness. If you reject them, you receive zero feedback about your decision's accuracy.

The patterns interviewers believe they recognize are often illusory, built from small sample sizes and confirmation bias. Yet in this environment of poor feedback and weak patterns, interviewers routinely express high confidence in their judgments.

# A Framework for Better Decisions
Kahneman suggests a framework.

**1. Define Key Traits**

Start by selecting 3-5 independent traits that are genuine prerequisites for role success. These might include technical proficiency, communication skills, problem-solving ability, or self-direction. These become your evaluation dimensions.

**2. Develop Specific Questions**

For each trait, create questions designed to evaluate that dimension alone. Avoid generic questions that blur multiple competencies.

**3. Pre-Commit to Scoring Criteria**

Before conducting interviews, establish clear scoring frameworks for each trait on a consistent scale (e.g., 1-5). Define what constitutes "very weak" versus "very strong" performance for each dimension. This pre-commitment prevents post-hoc rationalization and anchoring effects.

**4. Score Independently**

Evaluate each trait separately without forming overall impressions. This sequence matters critically when interviewers form holistic impressions first, those impressions contaminate specific competency evaluations.

**5. Synthesize with Constrained Intuition**

Only after systematic evaluation should experienced interviewers apply intuitive judgment to synthesize separate assessments. Gut feelings about cultural fit or leadership potential can meaningfully supplement structured evaluations—but only after objective work is complete.

Research in this area offers clear promise a structured approaches significantly outperform traditional unstructured interviews in low-validity environments.

> Note: This doesn't mean intuition is worthless in hiring—but it must be properly structured. Research consistently shows that even in low-validity environments, human judgment can add value when constrained by systematic processes.

# Alternative

**Technical Challenges**

Replace generic LeetCode problems with role specific technical scenarios. Instead of algorithmic puzzles, present candidates with realistic problems they'd encounter in the actual role. Provide code snippets with performance issues or architectural problems and ask them to identify and solve real-world challenges. Score based on predetermined criteria: problem identification (1-5), solution approach (1-5), and implementation details (1-5).

**PR Review**

PR reviews, for instance, are a core part of software engineering work, yet I've never been asked about my approach to code reviews. This is a missed opportunity because PR reviews easily distinguish different types of engineers revealing critical distinctions. 

Some candidates treat PR reviews merely as code standards enforcement and bug prevention, while others use them as platforms to connect, understand context and ideas, and stay in the loop with team decisions. These approaches offer fundamentally different signals about collaboration, mentorship ability, and technical leadership. Score independently on: technical issue identification (1-5), communication style (1-5), contextual understanding (1-5), and constructive feedback quality (1-5).

**System design alternative**

Standardize the evaluation process while keeping the open-ended nature. Pre-define the specific traits you're assessing (time management, solution completeness, tradeoff discussion, communication clarity). Create a scoring rubric for each trait before the interview. Train interviewers to evaluate these dimensions independently rather than looking for their personal preferred approach. Provide structured feedback forms that force interviewers to score each dimension separately before forming an overall impression.

**Behavioral Interviews** 

Structure the questions and scoring systematically and that means creating behavioral anchors for scoring (1-5 scale) that define what constitutes strong vs. weak responses for each competency. 
Notably, larger companies tend to do better at behavioral interviews precisely because they invest in developing these structured anchors and training interviewers to use them consistently, while smaller companies often rely on ad-hoc questioning with no standardized evaluation criteria.

The key change across all stages is to move from subjective impressions to objective, predetermined criteria applied consistently across all candidates.

It's worth acknowledging that hiring is inherently difficult and uncertain for both parties. Companies are trying to predict future performance based on limited interactions, while candidates are attempting to assess company culture, growth opportunities, and role fit through an artificial interview process. No framework can eliminate this fundamental uncertainty, but structured approaches can at least ensure that decisions are based on relevant signals rather than cognitive biases and subjective impressions.