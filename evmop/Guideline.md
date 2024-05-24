# Guideline

## Requirement
```bash
# create the project
mkdir evmop
cd evmop
operator-sdk init --domain gocrazy.com --repo evmop

# create an api
operator-sdk create api --group learn --version v1alpha1 --kind Blockchain --resource --controller
```