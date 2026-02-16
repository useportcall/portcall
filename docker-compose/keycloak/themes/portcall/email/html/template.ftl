<#macro emailLayout>
<!DOCTYPE html>
<html lang="${locale.language}" dir="${(ltr)?then('ltr','rtl')}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Portcall</title>
    <!--[if mso]>
    <style type="text/css">
        table, td { font-family: Segoe UI, Helvetica, Arial, sans-serif; }
    </style>
    <![endif]-->
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: #fafafa;
            font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
        }
        .email-wrapper {
            width: 100%;
            background-color: #fafafa;
            padding: 48px 0;
        }
        .email-container {
            max-width: 560px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            border: 1px solid #e5e5e5;
            overflow: hidden;
        }
        .email-header {
            padding: 32px 40px 0 40px;
        }
        .logo {
            font-size: 20px;
            font-weight: 600;
            color: #0a0a0a;
            text-decoration: none;
            letter-spacing: -0.4px;
        }
        .divider {
            height: 1px;
            background-color: #e5e5e5;
            margin: 24px 0 0 0;
            border: none;
        }
        .email-body {
            padding: 32px 40px 40px 40px;
        }
        .email-body h1 {
            color: #0a0a0a;
            font-size: 20px;
            font-weight: 600;
            margin: 0 0 8px 0;
            letter-spacing: -0.3px;
        }
        .email-body p {
            color: #737373;
            font-size: 14px;
            line-height: 1.6;
            margin: 0 0 16px 0;
        }
        .email-body a {
            color: #0a0a0a;
            text-decoration: underline;
            text-underline-offset: 2px;
        }
        .button {
            display: inline-block;
            background-color: #171717;
            color: #fafafa !important;
            padding: 10px 24px;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 500;
            text-decoration: none !important;
            margin: 8px 0;
            line-height: 1.5;
        }
        .link-box {
            word-break: break-all;
            background-color: #f5f5f5;
            padding: 12px 16px;
            border-radius: 8px;
            font-size: 13px;
            border: 1px solid #e5e5e5;
        }
        .link-box a {
            color: #737373;
            text-decoration: none;
        }
        .hint {
            color: #a3a3a3 !important;
            font-size: 13px !important;
            margin-top: 24px !important;
        }
        .email-footer {
            padding: 24px 40px;
            border-top: 1px solid #e5e5e5;
        }
        .email-footer p {
            color: #a3a3a3;
            font-size: 12px;
            margin: 0 0 4px 0;
            line-height: 1.5;
        }
        .email-footer a {
            color: #a3a3a3;
            text-decoration: underline;
            text-underline-offset: 2px;
        }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <div class="email-container">
            <div class="email-header">
                <span class="logo">Portcall</span>
                <hr class="divider">
            </div>
            <div class="email-body">
                <#nested>
            </div>
            <div class="email-footer">
                <p>&copy; ${.now?string('yyyy')} Portcall &middot; <a href="https://useportcall.com">useportcall.com</a></p>
            </div>
        </div>
    </div>
</body>
</html>
</#macro>
