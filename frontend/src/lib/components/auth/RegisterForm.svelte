<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import Input from '../ui/Input.svelte';
  import Button from '../ui/Button.svelte';
  import Alert from '../ui/Alert.svelte';

  export let loading = false;
  export let error: string | null = null;

  let name = '';
  let email = '';
  let password = '';
  let confirmPassword = '';
  let emailError = '';
  let passwordError = '';
  let confirmError = '';

  const dispatch = createEventDispatcher();

  function validateForm() {
    let isValid = true;
    emailError = '';
    passwordError = '';
    confirmError = '';

    if (!email) {
      emailError = 'Email is required';
      isValid = false;
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      emailError = 'Invalid email format';
      isValid = false;
    }

    if (!password) {
      passwordError = 'Password is required';
      isValid = false;
    } else if (password.length < 6) {
      passwordError = 'Password must be at least 6 characters';
      isValid = false;
    }

    if (password !== confirmPassword) {
      confirmError = 'Passwords do not match';
      isValid = false;
    }

    return isValid;
  }

  function handleSubmit() {
    if (!validateForm()) return;
    dispatch('submit', { email, password, name });
  }
</script>

<form on:submit|preventDefault={handleSubmit} class="space-y-4">
  {#if error}
    <Alert type="error" dismissible on:dismiss={() => error = null}>
      {error}
    </Alert>
  {/if}

  <Input
    label="Name (optional)"
    type="text"
    bind:value={name}
    placeholder="John Doe"
    disabled={loading}
  />

  <Input
    label="Email"
    type="email"
    bind:value={email}
    error={emailError}
    placeholder="you@example.com"
    disabled={loading}
  />

  <Input
    label="Password"
    type="password"
    bind:value={password}
    error={passwordError}
    placeholder="••••••"
    disabled={loading}
  />

  <Input
    label="Confirm Password"
    type="password"
    bind:value={confirmPassword}
    error={confirmError}
    placeholder="••••••"
    disabled={loading}
  />

  <Button type="submit" variant="primary" size="lg" disabled={loading} class="w-full">
    {loading ? 'Creating account...' : 'Register'}
  </Button>
</form>