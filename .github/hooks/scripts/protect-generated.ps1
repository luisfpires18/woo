$input = [Console]::In.ReadToEnd() | ConvertFrom-Json

$toolName = $input.toolName
$toolInput = $input.toolInput

# Check if this is a file edit tool targeting generated config files
$targetPath = $null
if ($toolInput.filePath) { $targetPath = $toolInput.filePath }
elseif ($toolInput.path) { $targetPath = $toolInput.path }

if ($targetPath -and $targetPath -match 'client[/\\]src[/\\]config[/\\]generated[/\\]') {
    $result = @{
        hookSpecificOutput = @{
            hookEventName = "PreToolUse"
            permissionDecision = "deny"
            permissionDecisionReason = "BLOCKED: Do not edit generated config files directly. Edit Go source in server/internal/config/ and run 'npm run gen-config' to regenerate."
        }
    }
    $result | ConvertTo-Json -Depth 5
    exit 0
}

# Allow all other operations
$result = @{
    hookSpecificOutput = @{
        hookEventName = "PreToolUse"
        permissionDecision = "allow"
    }
}
$result | ConvertTo-Json -Depth 5
exit 0
