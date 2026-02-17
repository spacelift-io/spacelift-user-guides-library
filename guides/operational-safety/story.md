# The Saturnhead Startup Chronicles: Operational Safety

_Orbit Labs is growing fast. Saturnhead's scrappy "move fast and break things" phase is over. Investors are asking about compliance. The new hire just pushed directly to main. It's time to add guardrails._

---

## Mission 1: Launchpad (Speedrun)

_"Before you can secure the mission, you need a mission."_

Saturnhead needs the basics in place — fast. VCS connected, AWS wired up, a stack that does something real.

### Tasks

1. Connect a VCS provider if you haven't already
2. Create a repository with Terraform that provisions an S3 bucket (we provide the code)
3. Create an AWS integration with cross-account role assumption
4. Create a stack, attach the integration
5. Trigger a run, confirm it completes

### API Checks

- VCS integration exists
- AWS integration exists with role assumption
- Stack exists with AWS integration attached
- At least one run in `FINISHED` state with resources created

---

## Mission 2: No Untagged Buckets

_"The finance team wants to know what's costing money. Saturnhead decrees: every resource gets a tag."_

Plan policies let you inspect _what_ Terraform is about to do — and stop it if it violates rules.

### Tasks

1. Create a plan policy requiring all S3 buckets to have a `cost-center` tag (we provide the OPA rule)
2. Attach it to your stack
3. Push a change that **violates** the policy — observe the run fail
4. Fix the Terraform — observe the run succeed

### API Checks

- Plan policy exists
- Policy attached to stack
- At least one run with policy `DENY`
- At least one subsequent run with policy `PASS`

### Failure Moment

Intentional violation → policy blocks it.

---

## Mission 3: Four Eyes

_"Saturnhead trusts his team. But production changes? Those need a second pair of eyes."_

Approval policies require humans to sign off before changes apply. Because some things shouldn't be one-person decisions.

### Tasks

1. Create an approval policy requiring one approval for any run (we provide the OPA rule)
2. Attach it to your stack
3. Trigger a run — observe it enters `UNCONFIRMED` state, waiting for approval
4. Approve the run (as yourself, for now)
5. Observe the run proceeds to apply

### API Checks

- Approval policy exists
- Policy attached to stack
- At least one run in state awaiting approval
- Run transitions to `FINISHED` after approval

### Learning Beat

The run _waits_. That's the point.

---

## Mission 4: Mission Control Knows

_"Saturnhead can't watch the dashboard 24/7. When something needs attention, he wants a ping."_

Notification policies route events to the right place. Start with the Inbox — Spacelift's built-in notification center.

### Tasks

1. Create a notification policy that sends to Inbox when a run needs approval (we provide the OPA rule)
2. Attach it to your stack
3. Trigger a run that requires approval
4. Check the Inbox — see the notification

### API Checks

- Notification policy exists (targeting inbox)
- Policy attached to stack
- Inbox contains notification related to the run

---

## Mission 5: Phone Home

_"Orbit Labs uses Slack. And incident.io. And a homegrown status board. Time to integrate."_

Webhooks let Spacelift talk to anything with an HTTP endpoint. Saturnhead sets up a test endpoint to see how it works.

### Tasks

1. Set up a webhook receiver (we provide: ngrok or webhook.site instructions)
2. Create a notification policy that sends run events to your webhook endpoint
3. Attach it to your stack
4. Trigger a run — observe the webhook payload arrive
5. Inspect the payload structure

### API Checks

- Notification policy exists in the space (targeting webhook)
- Policy attached to stack
- Run triggered
- Policy receipt present and correctly pointing at the created webhook
- Verify webhook receiver returned a response (most likely a 5xx since it's a test endpoint, but it should receive the payload and respond)
