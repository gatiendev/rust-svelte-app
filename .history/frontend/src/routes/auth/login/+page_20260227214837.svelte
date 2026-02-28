<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import LoginForm from '$lib/components/auth/LoginForm.svelte';
  import { auth } from '$lib/stores/auth';

  let loading = $auth.loading;
  let error = $auth.error;

  function handleLogin(event: CustomEvent<{ email: string; password: string }>) {
    const { email, password } = event.detail;
    auth.login(email, password);
  }
</script>

<div class="container mx-auto px-4 py-12 flex justify-center">
  <Card class="w-full max-w-md">
    <h1 class="text-2xl font-bold text-center mb-6">Login to Your Account</h1>
    <LoginForm on:submit={handleLogin} {loading} {error} />
    <p class="mt-4 text-center text-sm text-gray-600 dark:text-gray-400">
      Don't have an account?
      <a href="/auth/register" class="text-blue-600 dark:text-blue-400 hover:underline">Register</a>
    </p>
  </Card>
</div>