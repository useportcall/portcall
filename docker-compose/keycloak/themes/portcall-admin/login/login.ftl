<#-- Admin portal login page with purple accents -->
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Admin Portal - Portcall</title>
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
      background: linear-gradient(135deg, #f8f7ff 0%, #f0edff 100%);
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
      border: 1px solid #e9d5ff;
      border-radius: 18px;
      box-shadow: 0 4px 32px 0 rgba(139, 92, 246, 0.12);
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
      margin-bottom: 12px;
    }
    .admin-badge {
      background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
      color: #fff;
      padding: 6px 16px;
      border-radius: 20px;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      margin-bottom: 8px;
      box-shadow: 0 2px 8px rgba(139, 92, 246, 0.25);
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
    input[type="text"], input[type="password"] {
      width: 100%;
      padding: 10px 12px;
      border: 1.5px solid #e9d5ff;
      border-radius: 8px;
      background: #faf8ff;
      font-size: 0.98rem;
      color: #18181b;
      outline: none;
      transition: border 0.2s, box-shadow 0.2s;
      box-shadow: 0 1px 2px 0 rgba(139, 92, 246, 0.03);
    }
    input[type="text"]:focus, input[type="password"]:focus {
      border: 1.5px solid #8b5cf6;
      background: #fff;
      box-shadow: 0 2px 8px 0 rgba(139, 92, 246, 0.15);
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
      background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
      color: #fff;
      border: none;
      border-radius: 12px;
      padding: 14px 0;
      font-size: 18px;
      font-weight: 500;
      cursor: pointer;
      margin-top: 10px;
      letter-spacing: -0.2px;
      box-shadow: 0 4px 12px rgba(139, 92, 246, 0.25);
      transition: transform 0.2s, box-shadow 0.2s;
    }
    .login-btn:hover {
      transform: translateY(-1px);
      box-shadow: 0 6px 16px rgba(139, 92, 246, 0.35);
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
      color: #8b5cf6;
      text-decoration: underline;
      margin: 0 4px;
      transition: color 0.2s;
      font-weight: 400;
      font-size: 14px;
    }
    .login-links a:hover {
      color: #7c3aed;
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
    <div class="admin-badge">Admin Portal</div>
    <div class="login-title">Sign in to Portcall</div>
    <#if message?has_content>
      <div style="width: 100%; display: flex;
        flex-direction: column;">
        <div class="error-message">
          ${message.summary?no_esc}
        </div>
      </div>
    </#if>
    <form id="kc-form-login" class="form" action="${url.loginAction}" method="post">
      <div class="form-group">
        <label for="username">Username</label>
  <input tabindex="1" id="username" name="username" type="text" autofocus autocomplete="username" value="${login.username!''}" required />
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input tabindex="2" id="password" name="password" type="password" autocomplete="current-password" required />
      </div>
      <#if realm.rememberMe && !realm.password && !realm.otp && !realm.registrationEmailAsUsername>
        <div class="form-group">
          <label><input type="checkbox" id="rememberMe" name="rememberMe" <#if login.rememberMe??>checked</#if> /> Remember me</label>
        </div>
      </#if>
      <button class="login-btn" type="submit" tabindex="3">Log in</button>
    </form>
    <div class="login-links">
      <#if realm.resetPasswordAllowed>
        <a href="${url.loginResetCredentialsUrl}">Forgot password?</a>
      </#if>
      <#if realm.registrationAllowed>
        <span>|</span>
        <a href="${url.registrationUrl}">Sign up</a>
      </#if>
    </div>
  </div>
</body>
</html>


