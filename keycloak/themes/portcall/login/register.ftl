<#-- Minimal, clean, greyscale register page for Keycloak -->
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Create your Portcall account</title>
  <style>
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
    }
    .login-logo {
      margin-bottom: 18px;
    }
    .login-title {
      font-size: 1.35rem;
      font-weight: 700;
      margin-bottom: 18px;
      color: #18181b;
      letter-spacing: -0.5px;
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
    input[type="text"], input[type="email"], input[type="password"] {
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
    }
    input[type="text"]:focus, input[type="email"]:focus, input[type="password"]:focus {
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

    <div class="login-title">Create your Portcall account</div>

    <#-- Global (non-field) message from Keycloak -->
    <#if message?has_content>
      <div class="error-message">
        ${message.summary?no_esc}
      </div>
    </#if>

    <form id="kc-register-form" action="${url.registrationAction}" method="post" novalidate>

      <div class="form-group" style="flex-direction: row; gap: 16px;">
        <div style="flex: 1; display: flex; flex-direction: column;">
          <label for="firstName">First name</label>
          <input id="firstName" name="firstName" type="text" autocomplete="given-name" value="${(register.formData.firstName)!''}" required />
          <#if messagesPerField.existsError('firstName')>
            <div class="field-error">${messagesPerField.getFirstError('firstName')?no_esc}</div>
          </#if>
        </div>
        <div style="flex: 1; display: flex; flex-direction: column;">
          <label for="lastName">Last name</label>
          <input id="lastName" name="lastName" type="text" autocomplete="family-name" value="${(register.formData.lastName)!''}" required />
          <#if messagesPerField.existsError('lastName')>
            <div class="field-error">${messagesPerField.getFirstError('lastName')?no_esc}</div>
          </#if>
        </div>
      </div>

      <div class="form-group">
        <label for="email">Email</label>
        <input id="email" name="email" type="email" autocomplete="email" value="${(register.formData.email)!''}" required />
        <#if messagesPerField.existsError('email')>
          <div class="field-error">${messagesPerField.getFirstError('email')?no_esc}</div>
        </#if>
      </div>

      <#-- Only show username when realm is NOT using email as username -->
      <#if !realm.registrationEmailAsUsername>
        <div class="form-group">
          <label for="username">Username</label>
          <input id="username" name="username" type="text" autocomplete="username" value="${(register.formData.username)!''}" required />
          <#if messagesPerField.existsError('username')>
            <div class="field-error">${messagesPerField.getFirstError('username')?no_esc}</div>
          </#if>
        </div>
      </#if>

      <div class="form-group">
        <label for="password">Password</label>
        <input id="password" name="password" type="password" autocomplete="new-password" required />
        <#if messagesPerField.existsError('password')>
          <div class="field-error">${messagesPerField.getFirstError('password')?no_esc}</div>
        </#if>
      </div>

      <div class="form-group">
        <label for="password-confirm">Confirm password</label>
        <input id="password-confirm" name="password-confirm" type="password" autocomplete="new-password" required />
        <#if messagesPerField.existsError('password-confirm')>
          <div class="field-error">${messagesPerField.getFirstError('password-confirm')?no_esc}</div>
        </#if>
      </div>

      <button class="login-btn" type="submit">Create account</button>
    </form>

    <div class="login-links">
      <a href="${url.loginUrl}">Back to sign in</a>
    </div>
  </div>
</body>
</html>
