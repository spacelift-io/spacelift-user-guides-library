# The Saturnhead Startup Chronicles: Configuration Reuse

_Orbit Labs now has three teams and a dozen stacks. Saturnhead is tired of copy-pasting the same environment variables everywhere. Every stack needs the same proxy setup. Every stack needs the same tagging policy. There has to be a better way._

---

## Mission 1: Launchpad (Speedrun)

_"Before you can reuse config, you need config to reuse."_

Saturnhead needs the basics in place — fast. VCS connected, AWS wired up, a stack that does something real.

### Tasks

1. Connect a VCS provider
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

## Mission 2: Inline Everything

_"It works. But it's getting messy."_

Saturnhead needs his stack to know which environment it's deploying to, and to run a setup script before Terraform initializes. He adds these directly to the stack — the quick and dirty way.

### Tasks

1. Add an environment variable to your stack: `DEPLOY_ENV=production`
2. Add a `before_init` hook that echoes "Initializing for $DEPLOY_ENV" (we provide the script)
3. Trigger a run — observe the hook output in the logs
4. Verify the environment variable is accessible in Terraform (we provide: `output` that reads from env)

### API Checks

- Stack has environment variable `DEPLOY_ENV` configured
- Stack has `before_init` hook configured
- Run logs contain hook output
- Run outputs contain environment variable value

### Learning Beat

This works. But what happens when you have 10 stacks that all need the same setup?

---

## Mission 3: Extract and Reuse

_"Copy-paste is not a strategy."_

Saturnhead realizes the same config belongs to multiple stacks. Time to extract it into a Context — a reusable bundle of environment variables, files, and hooks.

### Tasks

1. Create a Context called `production-config`
2. Move the `DEPLOY_ENV` variable into the Context
3. Move the `before_init` hook into the Context
4. Remove the inline config from your stack
5. Attach the Context to your stack manually
6. Trigger a run — observe the same behavior as before

### API Checks

- Context exists with environment variable and hook
- Stack has context attached
- Stack no longer has inline environment variable or hook
- Run logs contain hook output
- Run outputs contain environment variable value

### Learning Beat

Same result. But now you can attach this Context to any stack that needs it.

---

## Mission 4: Labels and Autoattachment

_"Saturnhead doesn't want to manually attach contexts to every new stack. He wants it to just happen."_

Autoattachment lets you define rules: "any stack with label X gets context Y." No manual wiring required.

### Tasks

1. Add a label to your stack: `env:production`
2. Modify your Context to autoattach to stacks with label `env:production`
3. Detach the Context from your stack manually
4. Observe the Context is still attached — via autoattachment
5. Create a second stack with the same label — observe it gets the Context automatically

### API Checks

- Stack has label `env:production`
- Context has autoattachment rule for `env:production`
- Context attached to both stacks (via autoattachment, not manual)
- Second stack run shows context environment variable and hook in effect

### Learning Beat

Labels become contracts. "If you're production, you get production config."

---

## Mission 5: Config Flows Downhill

_"Orbit Labs has teams. Teams have environments. Some config applies to everything underneath."_

Saturnhead builds a space hierarchy. He wants config defined at the top to cascade down automatically.

### Tasks

1. Create a child space under root: `platform-team`
2. Create a grandchild space: `platform-team/production`
3. Move your stack into `platform-team/production`
4. Create a new Context with a company-wide setting: `COMPANY=orbit-labs`
5. Attach it at `platform-team` with `autoattach:*`
6. Trigger a run on your stack — observe the inherited environment variable in the logs

### API Checks

- Space hierarchy exists: `platform-team` → `platform-team/production`
- Stack located in `platform-team/production`
- Context exists with `COMPANY` environment variable
- Context attached at `platform-team` with `autoattach:*`
- Stack in child space has context attached (via inheritance)
- Run outputs contain `COMPANY=orbit-labs`

### Learning Beat

Define it once, inherit it everywhere. The space tree becomes your configuration tree.

---

## Mission 6: One Mechanism, Many Components

_"Wait — does this work for everything?"_

Saturnhead realizes autoattachment isn't just for Contexts. Policies and AWS integrations follow the same pattern. Define once, cascade everywhere.

### Tasks

1. Create a plan policy requiring a `team` tag on all resources (we provide the OPA rule)
2. Attach the policy at `platform-team` with `autoattach:*`
3. Detach the AWS integration from your stack manually
4. Attach the AWS integration at `platform-team` with `autoattach:*`
5. Trigger a run — observe:
   - The policy evaluates (and fails, since your bucket lacks the tag)
   - The AWS credentials are available (via inherited integration)
6. Fix the Terraform to add the tag — observe the run succeeds

### API Checks

- Policy attached at `platform-team` with `autoattach:*`
- AWS integration attached at `platform-team` with `autoattach:*`
- Stack in child space has both policy and integration attached (via inheritance)
- At least one run with policy `DENY`
- At least one subsequent run with policy `PASS`
- Run uses AWS credentials from inherited integration

### Learning Beat

Contexts, policies, integrations — same mechanism. The space hierarchy is your organizational model, and autoattachment is how you express "this applies to everyone below."
