# Percona DBaaS Platform Catalog

<!-- toc -->
- [Percona DBaaS Platform Catalog](#percona-dbaas-platform-catalog)
  - [Summary](#summary)
  - [Motivation](#motivation)
  - [Proposal](#proposal)
  - [Design Details](#design-details)
    - [Owners](#owners)
    - [Workflow](#workflow)
    - [Test Plan](#test-plan)
  - [Implementation History](#implementation-history)
<!-- /toc -->

## Summary

<!--
This section is incredibly important for producing high-quality, user-focused
documentation such as release notes or a development roadmap. It should be
possible to collect this information before implementation begins, in order to
avoid requiring implementors to split their attention between writing release
notes and implementing the feature itself. KEP editors and SIG Docs
should help to ensure that the tone and content of the `Summary` section is
useful for a wide audience.

A good summary is probably at least a paragraph in length.

Both in this section and below, follow the guidelines of the [documentation
style guide]. In particular, wrap lines to a reasonable length, to make it
easier for reviewers to cite specific portions, and to minimize diff churn on
updates.

[documentation style guide]: https://github.com/kubernetes/community/blob/master/contributors/guide/style-guide.md
-->

Percona DBaaS Platform OLM Catalog is a selection of operators that are tested and validated for Percona DBaaS products.

## Motivation

<!--
This section is for explicitly listing the motivation, goals, and non-goals of
this KEP.  Describe why the change is important and the benefits to users. The
motivation section can optionally provide links to [experience reports] to
demonstrate the interest in a KEP within the wider Kubernetes community.

[experience reports]: https://github.com/golang/go/wiki/ExperienceReports
-->

To create well tested and well maintained catalog to enable Percona platform in k8s.

## Proposal

<!--
This is where we get down to the specifics of what the proposal actually is.
This should have enough detail that reviewers can understand exactly what
you're proposing, but should not include things like API designs or
implementation. What is the desired outcome and how do we measure success?.
The "Design Details" section below is for the real
nitty-gritty.
-->

Should be an implementation of a platform as described in a [design document](dbaas-catalog-design.md#create-new-platform).

## Design Details

<!--
This section should contain enough information that the specifics of your
change are understandable. This may include API specs (though not always
required) or even code snippets. If there's any ambiguity about HOW your
proposal will be implemented, this is the place to discuss them.
-->

Create the branch `percona-platform` that would be as close as possible to a `main` but would have additional set of tests and might include/exclude some of the operators.

Operator's team creates PR (manually or automatically with CI) and changes first land to `main` and then merged to the `percona-platform`. Merge to `percona-platform` happens either automatically or manually by owners (after some quality gate).

Create `docker.io/percona/dbaas-catalog` registry to store catalog images. Bundles use same operator's name `docker.io/percona/percona-server-mongodb-operator:v1.13.0-bundle`) with `vX.Y.Z-bundle` tag. Testing and RC images are pushed to `percona-lab` registry.

Images would be created automatically by CI in `dbaas-catalog`:
  - with any commit to `percona-platform`: docker.io/percona/dbaas-catalog:main
  - when new tag pushed to `percona-platform` branch:
    - `docker.io/percona/dbaas-catalog:latest`
    - `docker.io/percona/dbaas-catalog:vX.Y.Z`

Images required from operators are bundle image:
  - **Candidate** Channel: `docker.io/percona-lab/operatorA:vX.Y.Z-rcN-bundle`
  - **Stable** Channel: `docker.io/percona/operatorA:vX.Y.Z-bundle`

### Owners

Platform owners are responsible and accountable for the platform CI/CD, quality (green actions) and state of the operators (new versions) in the catalog.

Platform owners:
  - @nmarukovich
  - @gen1us2k
  - @cap1984
  - @beatahandzelova

Owners of specific operators:
  - @nmarukovich
    - PSMDB
    - PXC
    - PS
  - @cap1984
    - PG
  - @gen1us2k
    - VM
    - DBaaS
  - @beatahandzelova
    - CI/CD
    - quality engineer for releases

### Workflow

Owners submit new versions of needed operators to the community catalog `main` and merge (or cherry-pick) needed changes to the `percona-platform` branch if/when quality gate is passing.

Owners are responsible for enhancing and maintaining CI/CD pipelines.

Owners create releases by adding tags to the `percona-platform` branch.

Bundles could come from the operator's registry or from other community catalogs/registries. But it is encouraged to have operators to build bundles continuously in their own registry and submit new versions (at least for **Candidates**) more frequently to enable more efficient interoperability testing.

### Test Plan

<!--
**Note:** *Not required until targeted at a release.*
The goal is to ensure that we don't accept enhancements with inadequate testing.

All code is expected to have adequate tests (eventually with coverage
expectations). Please adhere to the [Kubernetes testing guidelines][testing-guidelines]
when drafting this test plan.

[testing-guidelines]: https://git.k8s.io/community/contributors/devel/sig-testing/testing.md
-->

Testing happens on:
  - PR to merge changes to `main`
  - on `main` branch with new changes merged
  - on `percona-platform` branch when new changes merges

Each of those steps for testing could have different coverage:
  - sanity on PRs
  - integration, interoperability on `main`
  - integration, interoperability, e2e on `percona-platform`
  - specific platform test on `percona-platform` before release

Candidate channel is used for early interoperability testing of upcoming operators.

## Implementation History

<!--
Major milestones in the lifecycle of a KEP should be tracked in this section.
Major milestones might include:
- the `Summary` and `Motivation` sections being merged, signaling SIG acceptance
- the `Proposal` section being merged, signaling agreement on a proposed design
- the date implementation started
- the first Kubernetes release where an initial version of the KEP was available
- the version of Kubernetes where the KEP graduated to general availability
- when the KEP was retired or superseded
-->
