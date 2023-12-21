# Changelog

## v1.15.0-rc.0

Summary of Changes
------------------

**Minor Changes:**
* gateway-api: Update API version for Reference Grant (#29811, @sayboras)
* helm: Add missing SA automount configuration (#29511, @ayuspin)
* helm: Added support for existing Cilium SPIRE NS (#29032, @PhilipSchmid)
* helm: Allow setting resources for the agent init containers (#29610, @ayuspin)

**Bugfixes:**
* cilium-preflight: use the k8s node name instead of relying on hostname (#29809, @marseel)
* endpoint: fix panic in RunMetadataResolver due to send on closed channel (#29615, @mhofstetter)
* Fix bug where deleted nodes would reappear in the cilium_node_connectivity_* metrics (#29566, @christarazi)
* Fix cleanup of AWS-related leftover iptables chains (#29448, @giorio94)
* Fix missing NODE_ADD Hubble peer messages in some cases (#28226, @AwesomePatrol)
* Fix possible disruption of long running, cross-cluster, pod to node traffic on agent restart (#29613, @giorio94)
* Fix potential deadlock that results in stale authentication entries in Cilium (#29082, @meyskens)
* metrics: fix issue where logging err/warn metric is never updated. (#29201, @tommyp1ckles)
* The DNS proxy will now compute a UDP checksum over the IPv6 response packet and the pseudo-header. (#29493, @danehans)

**CI Changes:**
* ci datapath-verifier: add connectivity test (#29633, @mhofstetter)
* ci-ipsec-e2e: Misc refactor + more keys (#29592, @brb)
* ci-ipsec-upgrade: Add vxlan w/ no EP routes (#29653, @brb)
* ci-ipsec-{e2e,upgrade}: Use lvh-kind (#29514, @brb)
* ci/ipsec: Skip waiting for images when skipping upgrade/dowgrade (#29793, @qmonnet)
* ci: add documentation check to documentation workflow (#29684, @mhofstetter)
* ci: always use full matrix for scheduled cloud-provider workflows (#29694, @mhofstetter)
* ci: disable preemptible VM & GKE clusters on tests based on GKE (#29607, @mhofstetter)
* Define PUSH_TO_DOCKER_HUB environment variable (#29644, @michi-covalent)
* Fix collecting of verifier logs in ci-verifier (#29752, @lmb)
* Fix exporting results to gs bucket. (#29587, @marseel)
* gh/workflows: Bump CLI to v0.15.18 #29849 (@brb)
* gh/workflows: Drop rading /proc in case of failure (#29855, @brb)
* gh: e2e: test conformance & upgrade with 5.4 kernel and EgressGW (#29651, @julianwiedmann)
* gha: add step to ensure presence/absence of the AWS iptables chains (#29670, @giorio94)
* gha: enable IPv6 in clustermesh upgrade/downgrade workflow (#29675, @giorio94)
* gha: Migrate from MetalLB to L2LB (#28926, @sayboras)
* gha: sig-servicemesh owns Ingress or Gateway API related workflows (#29812, @sayboras)
* Make LB-IPAM tests less flaky (#29678, @dylandreimerink)
* Mock out time for BPF ratelimit test to make it more stable (#29740, @dylandreimerink)
* renovate: enable Cilium CLI patch updates for Cilium <v1.14 (#29794, @giorio94)
* Simplify CI image build workflow before v1.15 branch (#29834, @joestringer)
* test: Fail ginkgo tests on warnings (#29624, @pchaigno)
* workflows: Make the conn-disrupt test more sensitive (#29623, @pchaigno)

**Misc Changes:**
* Address device <-> node addressing race (#29555, @bimmlerd)
* bpf/Makefile: remove gen_compile_commands make target (#29611, @ti-mo)
* bpf: clean up some IPv4 header validations (#29585, @julianwiedmann)
* bpf: l3: restore MARK_MAGIC_PROXY_INGRESS for from-proxy traffic (#29721, @julianwiedmann)
* chore(deps): update actions/setup-python action to v4.8.0 (main) (#29769, @renovate[bot])
* chore(deps): update actions/stale action to v9 (main) (#29772, @renovate[bot])
* chore(deps): update all github action dependencies to v5 (main) (major) (#29773, @renovate[bot])
* chore(deps): update all lvh-images main (main) (patch) (#29556, @renovate[bot])
* chore(deps): update all lvh-images main (main) (patch) (#29766, @renovate[bot])
* chore(deps): update dependency cilium/cilium-cli to v0.15.17 (main) (#29557, @renovate[bot])
* chore(deps): update docker.io/library/alpine docker tag to v3.19.0 (main) (#29770, @renovate[bot])
* chore(deps): update docker.io/library/golang:1.21.5 docker digest to 2ff79bc (main) (#29765, @renovate[bot])
* chore(deps): update gcr.io/etcd-development/etcd docker tag to v3.5.11 (main) (#29767, @renovate[bot])
* chore(deps): update github/codeql-action action to v2.22.9 (main) (#29768, @renovate[bot])
* chore(deps): update go to v1.21.5 (main) (patch) (#29659, @renovate[bot])
* chore(deps): update google-github-actions/setup-gcloud action to v2 (main) (#29780, @renovate[bot])
* chore(deps): update hubble cli to v0.12.3 (main) (patch) (#29749, @renovate[bot])
* chore(deps): update quay.io/lvh-images/kind docker tag to bpf-next-20231211.012942 (main) (#29777, @renovate[bot])
* chore: add SI Analytics as cilium user (#29744, @JhoLee)
* chore: rename CIDRGroups resource to CiliumCIDRGroups (#29515, @pippolo84)
* cilium-dbg: Add "statedb node-addresses" command (#29479, @joamaki)
* cilium: Do not warn on socket tracing if EnableSocketLBTracing was not set (#29730, @borkmann)
* cilium: iptables masquerade to route source fixes (#29591, @borkmann)
* Clean up deprecated and unused IPCache APIs after FQDN transition to asynchronous APIs (#29657, @tklauser)
* CODEOWNERS: assign pkg/ip to @cilium/sig-agent (#29669, @tklauser)
* CODEOWNERS: sig-clustermesh additionally owns clustermesh-related GHA workflows and helm templates (#29671, @giorio94)
* codeowners: use new teams cilium/envoy & cilium/fqdn (#29627, @mhofstetter)
* daemon: Fix incorrect node and ciliumnode resource type in annotations (#29522, @hargrovee)
* do not start bandwidth manager in dry mode (#29183, @dylandreimerink)
* docs: add documentation for policy-cidr-match-mode=nodes (#28421, @squeed)
* docs: add MaxConnectedClusters documentation (#29637, @thorn3r)
* Docs: Adds Webhook Limitation to EKS Install Doc (#29497, @danehans)
* docs: Modify BGP MD5 password with Helm default change (#29527, @YutaroHayakawa)
* docs: specify which further release for fqdn option removal. (#29531, @squeed)
* Don't log an error if the to be deleted ipset entry does not exist (#29561, @giorio94)
* Envoy silence expected internal listener warning (#29786, @jrajahalme)
* envoy: perform version check directly on envoy binary (not starter) (#29512, @mhofstetter)
* examples: update guestbook example with new image registry (#29603, @mhofstetter)
* fix(deps): update all go dependencies main (main) (minor) (#29771, @renovate[bot])
* fix(deps): update all go dependencies main (main) (patch) (#29593, @renovate[bot])
* fqdn: avoid converting from `netip.Addr` to `net.IP` and back (#29625, @tklauser)
* guestbook: update example with leader/follower naming (#29642, @mhofstetter)
* helm: Allow unsupported K8s versions for now (#29888, @gandro)
* hubble-relay: fix panic during server shutdown (#29705, @mhofstetter)
* images: bump cni plugins to v1.4.0 (#29622, @squeed)
* improve the correctness of the rate limiting implementation in certain edge cases. (#29397, @dylandreimerink)
* ingress: add unit tests to test default ingressclass (#29792, @mhofstetter)
* ipcache: use TriggerController, not UpdateController (#29548, @squeed)
* k8s/resource: Add support for releasable `Resource[T]` (#29414, @pippolo84)
* Makefile: Fix variable override not working in all cases (#29599, @gandro)
* Optimize IP/FQDN management in the DNSCache (#29691, @squeed)
* pkg/rand: remove random  name generator (#29664, @aanm)
* pkg: proxy: only install from-proxy rules/routes for native routing (#29761, @julianwiedmann)
* plugins/cilium-cni: Introduce endpoint customization (#29707, @gandro)
* Prepare for release v1.15.0-pre.3 (#29596, @aanm)
* Prepare v1.15 stable branch (#29838, @joestringer)
* proxy: export ProxyConfig fields (#29827, @tklauser)
* README: Update releases (#29609, @aanm)
* release image: Allow arbitrary pre-release identifiers (#29173, @michi-covalent)
* Revert "cilium: Ensure xfrm state is initialized for route IP before … (#29801, @jrfastab)
* statedb: Fix revision indexing (#29840, @joamaki)
* test: remove probes-test.sh (#29612, @rgo3)
* Update SPIRE dependency to v1.8.5 (#29597, @meyskens)