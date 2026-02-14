<#-- Info/Confirmation page for Keycloak -->
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>${msg("loginTitle",(realm.displayName!''))}</title>
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
      text-align: center;
    }
    .info-message {
      font-size: 0.95rem;
      color: #52525b;
      text-align: center;
      line-height: 1.6;
      width: 100%;
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
    .warning-message {
      display: flex;
      color: #b45309;
      background: #fffbeb;
      border-radius: 7px;
      padding: 10px 14px;
      font-size: 1.01rem;
      margin-bottom: 18px;
      width: 100%;
      text-align: left;
      border: 1px solid #fde68a;
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
      text-decoration: none;
      text-align: center;
      display: block;
    }
    .login-btn:hover {
      background: #27272a;
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
  </style>
</head>
<body>
  <div class="login-container">
    <div class="login-logo">
      <img src="${url.resourcesPath}/img/logo.png" alt="Portcall logo" height="40"/>
    </div>
    
    <#if messageHeader??>
      <div class="login-title">${messageHeader}</div>
    <#else>
      <div class="login-title">${message.summary?no_esc}</div>
    </#if>
    
    <#if message?has_content && message.type != 'info'>
      <#if message.type = 'success'>
        <div class="success-message">${message.summary?no_esc}</div>
      <#elseif message.type = 'warning'>
        <div class="warning-message">${message.summary?no_esc}</div>
      <#else>
        <div class="error-message">${message.summary?no_esc}</div>
      </#if>
    <#elseif requiredActions??>
      <div class="info-message">
        <#list requiredActions as reqAction>
          ${msg("requiredAction.${reqAction}")}<#sep>, </#sep>
        </#list>
      </div>
    </#if>
    
    <#if skipLink??>
    <#else>
      <#if pageRedirectUri?has_content>
        <a class="login-btn" href="${pageRedirectUri}">Continue</a>
      <#elseif actionUri?has_content>
        <a class="login-btn" href="${actionUri}">Continue</a>
      <#elseif (client.baseUrl)?has_content>
        <a class="login-btn" href="${client.baseUrl}">Continue</a>
      </#if>
    </#if>
    
    <div class="login-links">
      <a href="https://dashboard.useportcall.com">
        Back to sign in
      </a>
    </div>
  </div>
</body>
</html>
