package com.tootest.ch.user.configuration

import io.jaegertracing.Configuration.*
import io.opentracing.Tracer
import org.slf4j.LoggerFactory
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration


@OnJaeger
@Configuration(proxyBeanMethods = false)
class JaegerConfiguration(private val jaegerProperties: JaegerProperties) {
    private val logger = LoggerFactory.getLogger(javaClass)

    @Bean
    fun tracer(): Tracer {
        logger.info("[JaegerProperties] {}", jaegerProperties)
        val config = io.jaegertracing.Configuration(jaegerProperties.serviceName)
        val sender = SenderConfiguration()
        sender.withEndpoint(jaegerProperties.endpoint)
        config.withSampler(SamplerConfiguration().withType("const").withParam(1))
        config.withReporter(ReporterConfiguration().withSender(sender).withMaxQueueSize(10000))
        return config.tracer
    }

}