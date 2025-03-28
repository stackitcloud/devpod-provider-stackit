name: stackit
version: {{ .Env.VERSION }}
description: |-
  DevPod on STACKIT
icon: https://github.com/stackitcloud/devpod-provider-stackit/blob/80d9a236cd309baf50e756480c4c00689f83b001/assets/stackit.png
optionGroups:
  - options:
      - AGENT_PATH
      - INACTIVITY_TIMEOUT
      - INJECT_DOCKER_CREDENTIALS
      - INJECT_GIT_CREDENTIALS
    name: "Agent options"
    defaultVisible: false
  - options:
    - STACKIT_PROJECT_ID
    - STACKIT_DISK_SIZE
    - STACKIT_FLAVOUR
    - STACKIT_REGION
    - STACKIT_PRIVATE_KEY_PATH
    - STACKIT_SERVICE_ACCOUNT_KEY_PATH
    - STACKIT_AVAILABILITY_ZONE
    name: "STACKIT options"
    defaultVisible: true
options:
  STACKIT_PROJECT_ID:
    description: "Project ID to use"
    required: true
    password: false
    default: ""
  STACKIT_DISK_SIZE:
    description: "The disk size to use."
    required: true
    password: false
    default: 64GB
  STACKIT_FLAVOR:
    description: "The machine type to use."
    default: g1.1
    suggestions:
    - g1.1
    - g1.2
    - g1.3
    - g1.4
    - g1a.1d
    - g1a.2d
    - g1a.4d
    - g1a.8d
    - g2i.1
    - g2i.2
    - g2i.4
    - g2i.8
    - b1.1
    - b1.2
    - b1a.1
    - b1a.2
    - b2i.1
    - b2i.2
    - c1.1
    - c1.2
    - c1.3
    - c1.4
    - c1.5
    - c1a.16
    - c1a.1
    - c1a.2
    - c1a.4
    - c1a.8
    - c2i.1
    - c2i.16
    - c2i.2
    - c2i.4
    - c2i.8
    - m1.1
    - m1.2
    - m1.3
    - m1a.1
    - m1a.2
    - m1a.4
    - m2i.1
    - m2i.2
    - m2i.4
    - s1.2
    - s1.3
    - s1.4
    - s1.5
    - s1a.16
    - s1a.2
    - s1a.4
    - s1a.8
    - t1.1
    - t1.2
    - t2i.1
  STACKIT_REGION:
    description: "The STACKIT region to create the VM in. e.g. eu01"
    required: true
    password: false
    suggestions:
    - eu01
    default: "eu01"
  STACKIT_AVAILABILITY_ZONE:
      description: "The STACKIT availability zone to create the VM in. e.g. eu01"
      required: true
      password: false
      suggestions:
      - eu01-1
      - eu01-2
      - eu01-3
      - eu01-M
      default: "eu01-1"
  STACKIT_PRIVATE_KEY_PATH:
    description: "The location of the custom private key file"
    required: false
    password: false
    default: ""
  STACKIT_SERVICE_ACCOUNT_KEY_PATH:
    description: "The location of the service account key file."
    required: true
    password: false
    default: ""
  INACTIVITY_TIMEOUT:
    description: If defined, will automatically stop the VM after the inactivity period.
    default: 10m
  INJECT_GIT_CREDENTIALS:
    description: "If DevPod should inject git credentials into the remote host."
    default: "true"
  INJECT_DOCKER_CREDENTIALS:
    description: "If DevPod should inject docker credentials into the remote host."
    default: "true"
  AGENT_PATH:
    description: The path where to inject the DevPod agent to.
    default: /var/lib/toolbox/devpod
agent:
  path: ${AGENT_PATH}
  dataPath: /home/devpod/.devpod
  inactivityTimeout: ${INACTIVITY_TIMEOUT}
  injectGitCredentials: ${INJECT_GIT_CREDENTIALS}
  injectDockerCredentials: ${INJECT_DOCKER_CREDENTIALS}
  binaries:
    STACKIT_PROVIDER:
{{- range file.Walk "./dist" -}}
{{- if not (file.IsDir .) -}}
    {{- $parts := . | regexp.Split "_" -1 }}
    {{- if eq "linux" (index $parts 1) }}
    - os: {{ index $parts 1 }}
      arch: {{ index $parts 2 }}
      path: https://github.com/stackitcloud/devpod-provider-stackit/releases/download/{{ $.Env.VERSION }}/{{ . | filepath.Base }}
      checksum: {{ . | file.Read | crypto.SHA256 }}
    {{- end -}}
{{- end -}}
{{- end }}
  exec:
    shutdown: |-
      ${STACKIT_PROVIDER} stop
binaries:
  STACKIT_PROVIDER:
{{- range file.Walk "./dist" -}}
{{- if not (file.IsDir .) -}}
    {{- $parts := . | regexp.Split "_" -1 }}
    {{- $ext := filepath.Ext . }}
    - os: {{ index $parts 1 }}
      arch: {{ index $parts 2 | strings.Trim $ext }}
      path: https://github.com/stackitcloud/devpod-provider-stackit/releases/download/{{ $.Env.VERSION }}/{{ . | filepath.Base }}
      checksum: {{ . | file.Read | crypto.SHA256 }}
{{- end -}}
{{- end }}
exec:
  init: ${STACKIT_PROVIDER} init
  command: ${STACKIT_PROVIDER} command
  create: ${STACKIT_PROVIDER} create
  delete: ${STACKIT_PROVIDER} delete
  start: ${STACKIT_PROVIDER} start
  stop: ${STACKIT_PROVIDER} stop
  status: ${STACKIT_PROVIDER} status