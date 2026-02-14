<#-- Password Reset Trigger - Enter email to receive reset link -->
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Reset your password</title>
  <style>
    @keyframes authFadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    @keyframes authCardIn {
      from { opacity: 0; transform: translateY(8px); }
      to { opacity: 1; transform: translateY(0); }
    }
    body {
      background: #fff;
      color: #222;
      font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
      margin: 0;
      padding: 0;
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      animation: authFadeIn 0.22s ease-out both;
    }
    .login-container {
      background: #fff;
      border: 1px solid #eee;
      border-radius: 18px;
      box-shadow: 0 4px 32px 0 rgba(0,0,0,0.07);
      padding: 56px 36px 40px 36px;
      width: 370px;
      max-width: 370px;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 18px;
      animation: authCardIn 0.26s ease-out both;
    }
    .login-logo {
      margin-bottom: 18px;
    }
    .login-title {
      font-size: 1.35rem;
      font-weight: 700;
      margin-bottom: 8px;
      color: #18181b;
      letter-spacing: -0.5px;
    }
    .login-subtitle {
      font-size: 0.95rem;
      color: #71717a;
      text-align: center;
      margin-bottom: 10px;
      line-height: 1.5;
    }
    .form-group {
      margin-bottom: 22px;
      display: flex;
      flex-direction: column;
      gap: 4px;
      width: 100%;
    }
    label {
      display: block;
      font-size: 0.93rem;
      color: #444;
      margin-bottom: 2px;
      font-weight: 500;
      letter-spacing: -0.2px;
    }
    input[type="text"], input[type="email"] {
      width: 100%;
      padding: 10px 12px;
      border: 1.5px solid #e5e7eb;
      border-radius: 8px;
      background: #f6f6f7;
      font-size: 0.98rem;
      color: #18181b;
      outline: none;
      transition: border 0.2s, box-shadow 0.2s;
      box-shadow: 0 1px 2px 0 rgba(0,0,0,0.01);
      box-sizing: border-box;
    }
    input[type="text"]:focus, input[type="email"]:focus {
      border: 1.5px solid #6366f1;
      background: #fff;
      box-shadow: 0 2px 8px 0 rgba(99,102,241,0.07);
    }
    .error-message {
      display: flex;
      color: #d32f2f;
      background: #f9eaea;
      border-radius: 7px;
      padding: 10px 14px;
      font-size: 1.01rem;
      margin-bottom: 18px;
      width: 100%;
      text-align: left;
      border: 1px solid #f2bdbd;
      box-sizing: border-box;
    }
    .success-message {
      display: flex;
      color: #15803d;
      background: #f0fdf4;
      border-radius: 7px;
      padding: 10px 14px;
      font-size: 1.01rem;
      margin-bottom: 18px;
      width: 100%;
      text-align: left;
      border: 1px solid #bbf7d0;
      box-sizing: border-box;
    }
    .login-btn {
      width: 100%;
      background: #18181b;
      color: #fff;
      border: none;
      border-radius: 12px;
      padding: 14px 0;
      font-size: 18px;
      font-weight: 500;
      cursor: pointer;
      margin-top: 10px;
      letter-spacing: -0.2px;
    }
    .login-links {
      margin-top: 22px;
      width: 100%;
      text-align: center;
      font-size: 1.01rem;
      color: #888;
      display: flex;
      flex-direction: row;
      justify-content: center;
      gap: 10px;
    }
    .login-links a {
      color: #6366f1;
      text-decoration: underline;
      margin: 0 4px;
      transition: color 0.2s;
      font-weight: 400;
      font-size: 14px;
    }
    .login-links a:hover {
      color: #18181b;
    }
    .form {
      width: 100%;
    }
  </style>
</head>
<body>
  <div class="login-container">
    <div class="login-logo">
      <img src="${url.resourcesPath}/img/logo.png" alt="Portcall logo" height="40"/>
    </div>
    <div class="login-title">Reset your password</div>
    <div class="login-subtitle">Enter your email address and we'll send you a link to reset your password.</div>
    <#if message?has_content>
      <div style="width: 100%; display: flex; flex-direction: column;">
        <#if message.type = 'success'>
          <div class="success-message">
            ${message.summary?no_esc}
          </div>
        <#else>
          <div class="error-message">
            ${message.summary?no_esc}
          </div>
        </#if>
      </div>
    </#if>
    <form id="kc-reset-password-form" class="form" action="${url.loginAction}" method="post">
      <div class="form-group">
        <label for="username">Email</label>
        <input tabindex="1" id="username" name="username" type="text" autofocus autocomplete="email" value="<#if auth?has_content && auth.showUsername()>${auth.attemptedUsername}</#if>" required />
      </div>
      <button class="login-btn" type="submit" tabindex="2">Send reset link</button>
    </form>
    <div class="login-links">
      <a href="${url.loginUrl}">Back to sign in</a>
    </div>
  </div>
</body>
</html>
