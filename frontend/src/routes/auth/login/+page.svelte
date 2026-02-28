<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import LoginForm from "$lib/components/auth/LoginForm.svelte";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";

  $: loading = $auth.loading;
  $: error = $auth.error;

  function handleLogin(
    event: CustomEvent<{ email: string; password: string }>,
  ) {
    const { email, password } = event.detail;
    auth.login(email, password);
  }

  // Redirect to home after successful login
  $: if ($auth.user) {
    goto("/");
  }
</script>

<div
  class="container mx-auto px-4 py-16 flex justify-center items-center min-h-[calc(100vh-200px)]"
>
  <div class="w-full max-w-md">
    <Card class="p-8 shadow-xl border border-gray-200 dark:border-gray-700">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          Welcome Back
        </h1>
        <p class="text-gray-600 dark:text-gray-400 mt-2">
          Sign in to your account
        </p>
      </div>
      <LoginForm on:submit={handleLogin} {loading} {error} />
      <p class="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
        Don't have an account?
        <a
          href="/auth/register"
          class="text-primary-600 dark:text-primary-400 hover:underline font-medium"
          >Create one</a
        >
      </p>
    </Card>
  </div>
</div>
