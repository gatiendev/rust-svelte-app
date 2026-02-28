<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import RegisterForm from '$lib/components/auth/RegisterForm.svelte';
  import { auth } from '$lib/stores/auth';

  let loading = $auth.loading;
  let error = $auth.error;

  function handleRegister(event: CustomEvent<{ email: string; password: string; name?: string }>) {
    const { email, password, name } = event.detail;
    auth.register(email, password, name);
  }
</script>

<div class="container mx-auto px-4 py-12 flex justify-center">
  <Card class="w-full max-w-md">
    <h1 class="text-2xl font-bold text-center mb-6">Create an Account</h1>
    <RegisterForm on:submit={handleRegister} {loading} {error} />
    <p class="mt-4 text-center text-sm text-gray-600 dark:text-gray-400">
      Already have an account?
      <a href="/auth/login" class="text-blue-600 dark:text-blue-400 hover:underline">Login</a>
    </p>
  </Card>
</div>