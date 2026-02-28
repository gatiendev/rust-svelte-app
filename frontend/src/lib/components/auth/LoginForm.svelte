<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import Input from '../ui/Input.svelte';
  import Button from '../ui/Button.svelte';
  import Alert from '../ui/Alert.svelte';

  export let loading = false;
  export let error: string | null = null;

  let email = '';
  let password = '';
  let emailError = '';
  let passwordError = '';

  const dispatch = createEventDispatcher();

  function validateForm() {
    let isValid = true;
    emailError = '';
    passwordError = '';

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

    return isValid;
  }

  function handleSubmit() {
    if (!validateForm()) return;
    dispatch('submit', { email, password });
  }
</script>

<form on:submit|preventDefault={handleSubmit} class="space-y-4">
  {#if error}
    <Alert type="error" dismissible on:dismiss={() => error = null}>
      {error}
    </Alert>
  {/if}

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

  <Button type="submit" variant="primary" size="lg" disabled={loading} class="w-full">
    {loading ? 'Logging in...' : 'Login'}
  </Button>
</form>