# UI/CLI Copy Language

This file is the source of truth for user-facing CLI copy language in VocabMaster. Update this file before changing the active copy language mode, default locale, or supported locales.

## Mode

`single-language`: Simplified Chinese (zh-CN)

## Scope

This file applies to text shown to CLI users: command output, command help, prompts, confirmations, validation messages, progress text, empty states, and user-facing terminal errors.

Learning content may include English and Japanese terms, readings, example sentences, dictionary fields, and mnemonics. Explanations around that content should be Chinese unless the lesson content intentionally requires otherwise.

This file does not govern code identifiers, protocol names, logs, developer comments, test names, third-party strings, commands, flags, provider names, file paths, data field names, config keys, or other technical identifiers unless those strings are rendered to users.

## Rules

- Write new user-facing CLI copy in Simplified Chinese (zh-CN).
- Keep commands, flags, provider names, file paths, data field names, config keys, and technical identifiers as written in code.
- If a user requests CLI copy in a language or locale not listed here, update this file first and then make the copy change.
- If project code and this file disagree, stop and update this file or ask for the intended language mode before changing user-facing CLI copy.
